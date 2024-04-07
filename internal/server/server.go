package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type dompserver struct {
	coreMux *chi.Mux
	coreStg storage.StorageInterface
}

func Start() {
	cfg := configs.CreateServerConfig()
	logger.Initialize(cfg.Loglevel, "server_")
	dompserv := CreateNewServer()
	if dompserv.coreMux == nil || dompserv.coreStg == nil {
		logger.Log.Info("Server initialization failed!", zap.Bool("coreMux", (dompserv.coreMux == nil)), zap.Bool("coreStg", (dompserv.coreStg == nil)))
		return
	}
	logger.Log.Info("Server was successfully initialized!")
	err := http.ListenAndServe(cfg.Address, dompserv.coreMux)
	if err != nil {
		panic(err)
	}
}

func CreateNewServer() *dompserver {
	coreMux := chi.NewRouter()
	coreStg := &storage.MemStorage{}
	coreStg.InitMemStorage()
	dompserv := &dompserver{coreMux: coreMux, coreStg: coreStg}

	coreMux.Use(dompserv.WithResponseLog)
	coreMux.Use(dompserv.WithRequestLog)
	coreMux.Use(gzipHandle)
	coreMux.Use(DecompressHandler)

	coreMux.Get("/", dompserv.GetMainPageHandler)
	coreMux.Route("/update", func(r chi.Router) {
		r.Post("/", dompserv.MetricHandlerJSON)
		r.Get("/", dompserv.IncorrectRequestHandler)
		r.Post("/{mType}/{mName}/{mValue}", dompserv.UpdateMetricHandler)
		r.Get("/{mType}/{mName}/{mValue}", dompserv.IncorrectRequestHandler)
	})
	coreMux.Route("/value", func(r chi.Router) {
		r.Post("/", dompserv.MetricHandlerJSON)
		r.Get("/{mType}/{mName}", dompserv.GetMetricHandler)
	})
	return dompserv
}
