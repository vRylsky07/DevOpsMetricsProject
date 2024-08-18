package server

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/backups/dompdb"
	filesbackup "DevOpsMetricsProject/internal/backups/files"
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type dompserver struct {
	coreMux *chi.Mux
	coreStg storage.MetricsRepository
	cfg     *configs.ServerConfig
	log     logger.Recorder
}

func (serv *dompserver) IsValid() bool {
	return serv.coreMux != nil || serv.coreStg.IsValid() || serv.cfg != nil || serv.log != nil
}

func Start() {
	dompserv, err := CreateNewServer(configs.CreateServerConfig())

	if err != nil {
		panic(err)
	}

	dompserv.log.Info("Server was successfully initialized!")

	err = http.ListenAndServe(dompserv.cfg.Address, dompserv.coreMux)

	if err != nil {
		panic(err)
	}

}

func CreateNewServer(cfg *configs.ServerConfig) (*dompserver, error) {
	dompserv, err := NewDompServer(cfg)

	if !dompserv.IsValid() {
		return nil, err
	}

	dompserv.coreMux.Use(dompserv.WithResponseLog)
	dompserv.coreMux.Use(dompserv.WithRequestLog)
	dompserv.coreMux.Use(dompserv.gzipHandle)
	dompserv.coreMux.Use(dompserv.DecompressHandler)

	dompserv.coreMux.Get("/", dompserv.GetMainPageHandler)
	dompserv.coreMux.Get("/ping", dompserv.PingDatabaseHandler)
	dompserv.coreMux.Route("/update", func(r chi.Router) {
		r.Post("/", dompserv.MetricHandlerJSON)
		r.Get("/", dompserv.IncorrectRequestHandler)
		r.Post("/{mType}/{mName}/{mValue}", dompserv.UpdateMetricHandler)
		r.Get("/{mType}/{mName}/{mValue}", dompserv.IncorrectRequestHandler)
	})
	dompserv.coreMux.Route("/value", func(r chi.Router) {
		r.Post("/", dompserv.MetricHandlerJSON)
		r.Get("/{mType}/{mName}", dompserv.GetMetricHandler)
	})
	dompserv.coreMux.Route("/updates", func(r chi.Router) {
		r.Post("/", dompserv.UpdateBatchHandler)
	})
	return dompserv, nil
}

func NewDompServer(cfg *configs.ServerConfig) (*dompserver, error) {
	var errs []error

	logger, errLog := logger.Initialize(cfg.Loglevel, "server_")

	if errLog != nil {
		errs = append(errs, errLog)
	}

	mode := ""
	switch cfg.SaveMode {
	case constants.DatabaseMode:
		mode = "IS DATABASE MODE"
	case constants.FileMode:
		mode = "IS FILE MODE"
	case constants.InMemoryMode:
		mode = "IS MEMORY MODE"
	}
	logger.Info(mode)
	coreMux := chi.NewRouter()

	var coreStg storage.MetricsRepository
	var backup backup.MetricsBackup
	var bckErr error

	switch cfg.SaveMode {
	case constants.DatabaseMode:
		backup, bckErr = dompdb.NewDompDB(cfg.DatabaseDSN, logger)
		if bckErr != nil {
			errs = append(errs, bckErr)
		}
		coreStg = storage.NewBackupSupportStorage(cfg.RestoreBool, backup, logger)
	case constants.FileMode:
		backup, bckErr = filesbackup.NewMetricsBackup(cfg, logger)
		if bckErr != nil {
			errs = append(errs, bckErr)
		}
		coreStg = storage.NewBackupSupportStorage(cfg.RestoreBool, backup, logger)
	case constants.InMemoryMode:
		coreStg = storage.NewMemStorage()
	}

	serv := &dompserver{
		coreMux: coreMux,
		coreStg: coreStg,
		cfg:     cfg,
		log:     logger,
	}

	return serv, errors.Join(errs...)
}
