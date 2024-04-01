package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/coretypes"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
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

	coreMux.Use(WithResponseLog)
	coreMux.Use(WithRequestLog)

	coreMux.Get("/", dompserv.GetMainPageHandler)
	coreMux.Route("/update", func(r chi.Router) {
		r.Post("/", dompserv.UpdateMetricHandlerJSON)
		r.Get("/", dompserv.IncorrectRequestHandler)
		r.Post("/{mType}/{mName}/{mValue}", dompserv.UpdateMetricHandler)
		r.Get("/{mType}/{mName}/{mValue}", dompserv.IncorrectRequestHandler)
	})
	coreMux.Route("/value", func(r chi.Router) {
		r.Get("/{mType}/{mName}", dompserv.GetMetricHandler)
	})
	return dompserv
}

func (serv *dompserver) GetMainPageHandler(res http.ResponseWriter, req *http.Request) {

	if serv == nil || serv.coreMux == nil || serv.coreStg == nil {
		http.Error(res, "Server not working fine please check its initialization!", http.StatusBadRequest)
		return
	}

	htmlTopPart := `<html>
    <head>
    <title></title>
    </head>
    <body>`

	htmlBottomPart := `</body>
	</html>`

	htmlMiddlePart := ``

	g, c := serv.coreStg.ReadMemStorageFields()

	gSortedNames := maps.Keys(g)
	sort.Slice(gSortedNames, func(i, j int) bool {
		return gSortedNames[i] < gSortedNames[j]
	})

	cSortedNames := maps.Keys(c)
	sort.Slice(cSortedNames, func(i, j int) bool {
		return cSortedNames[i] < cSortedNames[j]
	})

	for _, key := range gSortedNames {
		value, errBool := g[key]
		if errBool {
			htmlMiddlePart += key + " " + strconv.FormatFloat(value, 'f', -1, 64) + "<br>"
		}
	}

	for _, key := range cSortedNames {
		value, errBool := c[key]
		if errBool {
			htmlMiddlePart += key + " " + strconv.Itoa(value) + "<br>"
		}
	}

	if htmlMiddlePart == "" {
		htmlMiddlePart = "SERVER STORAGE IS EMPTY FOR NOW"
	}

	htmlFinal := htmlTopPart + htmlMiddlePart + htmlBottomPart

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlFinal))
}

func (serv *dompserver) GetMetricHandler(res http.ResponseWriter, req *http.Request) {
	mType := chi.URLParam(req, "mType")
	mName := chi.URLParam(req, "mName")

	if mType != "gauge" && mType != "counter" {
		http.Error(res, "Your request is incorrect! Metric type should be `gauge` or `counter`", http.StatusBadRequest)
		return
	}

	if serv == nil {
		http.Error(res, "ERROR! Server not working fine please check its initialization! Server is nil", http.StatusBadRequest)
		return
	}

	if serv.coreMux == nil {
		http.Error(res, "ERROR! Server not working fine please check its initialization! Serv.coreMux is nil", http.StatusBadRequest)
		return
	}

	if serv.coreStg == nil {
		http.Error(res, "ERROR! Server not working fine please check its initialization! Serv.coreStg is nil", http.StatusBadRequest)
		return
	}

	mTypeConst := functionslibrary.ConvertStringToMetricType(mType)

	if mTypeConst == constants.NoneType {
		http.Error(res, "Your request is incorrect! Metric type conversion failed!", http.StatusBadRequest)
		return
	}

	valueToReturn, gettingValueError := serv.coreStg.GetMetricByName(mTypeConst, mName)

	if gettingValueError == nil {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatFloat(valueToReturn, 'f', -1, 64)))
		return
	} else {
		http.Error(res, "This metric does not exist or was not been updated yet", http.StatusNotFound)
		return
	}
}

func (serv *dompserver) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {

	mType := chi.URLParam(req, "mType")
	mName := chi.URLParam(req, "mName")
	mValue := chi.URLParam(req, "mValue")

	if mName == "" {
		http.Error(res, "Your request is incorrect! Please enter valid metric name", http.StatusNotFound)
		return
	}

	mTypeConst := functionslibrary.ConvertStringToMetricType(mType)

	valueInFloat, err := strconv.ParseFloat(mValue, 64)

	if (mType == "gauge" || mType == "counter") && err == nil {

		switch mTypeConst {
		case constants.GaugeType:
			serv.coreStg.UpdateMetricByName(constants.RenewOperation, mTypeConst, mName, valueInFloat)
		case constants.CounterType:
			serv.coreStg.UpdateMetricByName(constants.AddOperation, mTypeConst, mName, valueInFloat)
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Metrics was been updated! Thank you!"))
		return
	}

	http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
}

func (serv *dompserver) IncorrectRequestHandler(res http.ResponseWriter, req *http.Request) {
	http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
}

func (serv *dompserver) UpdateMetricHandlerJSON(res http.ResponseWriter, req *http.Request) {

	var mReceiver coretypes.Metrics
	err := json.NewDecoder(req.Body).Decode(&mReceiver)

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	switch functionslibrary.ConvertStringToMetricType(mReceiver.MType) {
	case constants.GaugeType:
		serv.coreStg.UpdateMetricByName(constants.RenewOperation, constants.GaugeType, mReceiver.ID, *mReceiver.Value)
		*mReceiver.Value, _ = serv.coreStg.GetMetricByName(constants.GaugeType, mReceiver.ID)
		logger.Log.Info(`Metric "` + mReceiver.ID + `" was successfully updated! New value is ` + strconv.FormatFloat(*mReceiver.Value, 'f', -1, 64))
		return

	case constants.CounterType:
		serv.coreStg.UpdateMetricByName(constants.AddOperation, constants.CounterType, mReceiver.ID, float64(*mReceiver.Delta))
		var counterValue float64
		counterValue, _ = serv.coreStg.GetMetricByName(constants.GaugeType, mReceiver.ID)
		*mReceiver.Delta = int64(counterValue)
		logger.Log.Info(`Metric "` + mReceiver.ID + `" was successfully updated! New value is ` + strconv.Itoa(int(*mReceiver.Delta)))
		return

	default:
		logger.Log.Error("ConvertStringToMetricType returns NoneType")
		return
	}

}

func WithRequestLog(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)

		logger.Log.Info("Server got HTTP-request", zap.String("uri", uri), zap.String("method", method), zap.Duration("time", duration))

	}
	return http.HandlerFunc(logFn)
}

func WithResponseLog(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		rlw := &ResponceLogWriter{w, 0, 0}
		h.ServeHTTP(rlw, r)

		logger.Log.Info("Server responsing", zap.Int("status", rlw.statusCode), zap.Int("size", rlw.size))
	}

	return http.HandlerFunc(logFn)
}

type ResponceLogWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rlw *ResponceLogWriter) Write(b []byte) (int, error) {
	size, err := rlw.ResponseWriter.Write(b)
	rlw.size = size
	return size, err
}

func (rlw *ResponceLogWriter) WriteHeader(statusCode int) {
	rlw.ResponseWriter.WriteHeader(statusCode)
	rlw.statusCode = statusCode
}
