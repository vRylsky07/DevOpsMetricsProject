package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type dompserver struct {
	coreMux        *chi.Mux
	coreStg        storage.StorageInterface
	currentMetrics *os.File
	cfg            *configs.ServerConfig
	savefile       *MetricsSave
	dompdb         *dompdb
}

func (serv *dompserver) IsValid() bool {
	if serv.coreMux != nil || serv.coreStg != nil || serv.currentMetrics != nil || serv.cfg != nil || serv.savefile != nil || serv.savefile.IsValid() {
		return true
	}
	logger.Log.Error("DOMP Server is not valid")
	return false
}

func Start() {
	dompserv := CreateNewServer(configs.CreateServerConfig())

	switch dompserv.cfg.SaveMode {
	case constants.FileMode:
		dompserv.StartSaveMetricsThread()
	case constants.DatabaseMode:
		if !dompserv.dompdb.IsValid() {
			return
		}
	}

	if !dompserv.IsValid() {
		logger.Log.Info(
			"Server initialization failed!  ",
			zap.Bool("coreMux is nil?", (dompserv.coreMux == nil)),
			zap.Bool("coreStg is nil?", (dompserv.coreStg == nil)),
			zap.Bool("currentMetrics is nil?", (dompserv.currentMetrics == nil)),
			zap.Bool("cfg is nil?", (dompserv.cfg == nil)),
			zap.Bool("savefile is not valid?", (!dompserv.savefile.IsValid())),
		)
		return
	}
	logger.Log.Info("Server was successfully initialized!")

	err := http.ListenAndServe(dompserv.cfg.Address, dompserv.coreMux)

	if err != nil {
		panic(err)
	}

}

func CreateNewServer(cfg *configs.ServerConfig) *dompserver {
	dompserv := NewDompServer(cfg)

	dompserv.coreMux.Use(dompserv.WithResponseLog)
	dompserv.coreMux.Use(dompserv.WithRequestLog)
	dompserv.coreMux.Use(gzipHandle)
	dompserv.coreMux.Use(DecompressHandler)

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
	return dompserv
}

func NewDompServer(cfg *configs.ServerConfig) *dompserver {
	coreMux := chi.NewRouter()
	coreStg := &storage.MemStorage{}
	coreStg.InitMemStorage()
	logger.Initialize(cfg.Loglevel, "server_")

	var currentMetrics *os.File = nil
	var db *dompdb

	switch cfg.SaveMode {
	case constants.DatabaseMode:
		var err error
		db, err = RunDB(cfg.DatabaseDSN)

		if err != nil {
			return nil
		}

	case constants.FileMode:
		currentMetrics = CreateTempFile(cfg.TempFile, cfg.RestoreBool)
	}

	serv := &dompserver{
		coreMux:        coreMux,
		coreStg:        coreStg,
		currentMetrics: currentMetrics,
		cfg:            cfg,
		savefile:       RestoreData(cfg, coreStg),
		dompdb:         db,
	}
	return serv
}
