package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type dompserver struct {
	coreMux        *chi.Mux
	coreStg        storage.MetricsRepository
	currentMetrics *os.File
	cfg            *configs.ServerConfig
	savefile       *MetricsSave
	db             storage.DompInterfaceDB
	log            logger.Recorder
}

func (serv *dompserver) IsValid() bool {
	if serv.coreMux != nil || serv.coreStg != nil || serv.currentMetrics != nil || serv.cfg != nil || serv.savefile != nil || serv.savefile.IsValid() || serv.log != nil {
		if (serv.cfg.SaveMode == constants.DatabaseMode) && (serv.db == nil || !serv.db.IsValid()) {
			return false
		}
		return true
	}
	return false
}

func Start() {
	dompserv, err := CreateNewServer(configs.CreateServerConfig())

	if err != nil {
		panic(err)
	}

	if dompserv.cfg.SaveMode == constants.FileMode {
		dompserv.StartSaveMetricsThread()
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
	coreMux := chi.NewRouter()
	coreStg := &storage.MemStorage{}
	coreStg.InitMemStorage()

	var errs []error

	logger, errLog := logger.Initialize(cfg.Loglevel, "server_")

	if errLog != nil {
		errs = append(errs, errLog)
	}

	var currentMetrics *os.File
	var db DompInterfaceDB

	switch cfg.SaveMode {
	case constants.DatabaseMode:
		var err error
		db, err = RunDB(cfg.DatabaseDSN, logger)

		if err != nil {
			errs = append(errs, errLog)
		}

	case constants.FileMode:
		currentMetrics = CreateTempFile(cfg.TempFile, cfg.RestoreBool, logger)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	serv := &dompserver{
		coreMux:        coreMux,
		coreStg:        coreStg,
		currentMetrics: currentMetrics,
		cfg:            cfg,
		savefile:       RestoreData(cfg, db, coreStg, logger),
		db:             db,
		log:            logger,
	}

	return serv, errors.Join(errs...)
}
