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
	pinger  backup.PingerDB
}

func (serv *dompserver) IsValid() bool {
	b := serv.coreMux != nil || serv.coreStg.IsValid() || serv.cfg != nil || serv.log != nil

	if serv.cfg.SaveMode == constants.DatabaseMode {
		return b && serv.pinger != nil
	}

	return b
}

func Start() {
	dompserv, err := CreateNewServer(configs.CreateServerConfig())

	if err != nil {
		panic(err)
	}

	dompserv.log.Info("Server was successfully initialized!")

	dompserv.log.Info("Starting server on address: " + dompserv.cfg.Address)
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
	dompserv.coreMux.Use(dompserv.GzipMiddleware)
	dompserv.coreMux.Use(dompserv.DecompressMiddleware)
	dompserv.coreMux.Use(dompserv.HashCompareMiddleware)

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

	switch cfg.SaveMode {
	case constants.DatabaseMode:
		logger.Info("Backup save mode: Database")
	case constants.FileMode:
		logger.Info("Backup save mode: File")
	case constants.InMemoryMode:
		logger.Info("Backup save mode: In memory")
	}

	coreMux := chi.NewRouter()
	var coreStg storage.MetricsRepository
	var backuper backup.MetricsBackup
	var bckErr error
	var pinger backup.PingerDB

	switch cfg.SaveMode {
	case constants.DatabaseMode:
		backuper, bckErr = dompdb.NewDompDB(cfg.DatabaseDSN, logger)
		if bckErr != nil {
			errs = append(errs, bckErr)
		}
		pinger = backuper.(backup.PingerDB)
		coreStg = storage.NewBackupSupportStorage(cfg.RestoreBool, backuper, logger)
	case constants.FileMode:
		backuper, bckErr = filesbackup.NewMetricsBackup(cfg, logger)
		if bckErr != nil {
			errs = append(errs, bckErr)
		}
		coreStg = storage.NewBackupSupportStorage(cfg.RestoreBool, backuper, logger)
	case constants.InMemoryMode:
		coreStg = storage.NewMemStorage()
	}

	serv := &dompserver{
		coreMux: coreMux,
		coreStg: coreStg,
		cfg:     cfg,
		log:     logger,
		pinger:  pinger,
	}

	return serv, errors.Join(errs...)
}
