package server

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/coretypes"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/logger"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

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
		logger.Log.ErrorHTTP(res, errors.New("your request is incorrect! Metric type should be `gauge` or `counter`"), http.StatusBadRequest)
		return
	}

	if serv == nil {
		logger.Log.ErrorHTTP(res, errors.New("ERROR! Server not working fine please check its initialization! Server is nil"), http.StatusBadRequest)
		return
	}

	if serv.coreMux == nil {
		logger.Log.ErrorHTTP(res, errors.New("ERROR! Server not working fine please check its initialization! Serv.coreMux is nil"), http.StatusBadRequest)
		return
	}

	if serv.coreStg == nil {
		logger.Log.ErrorHTTP(res, errors.New("ERROR! Server not working fine please check its initialization! Serv.coreStg is nil"), http.StatusBadRequest)
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

	if serv == nil {
		logger.Log.ErrorHTTP(res, errors.New("ERROR! Server not working fine please check its initialization! Server is nil"), http.StatusBadRequest)
		return
	}

	if serv.coreMux == nil {
		logger.Log.ErrorHTTP(res, errors.New("ERROR! Server not working fine please check its initialization! Serv.coreMux is nil"), http.StatusBadRequest)
		return
	}

	if serv.coreStg == nil {
		logger.Log.ErrorHTTP(res, errors.New("ERROR! Server not working fine please check its initialization! Serv.coreStg is nil"), http.StatusBadRequest)
		return
	}

	var mReceiver coretypes.Metrics
	err := json.NewDecoder(req.Body).Decode(&mReceiver)

	if err != nil {
		logger.Log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	var newValue float64

	stringType := functionslibrary.ConvertStringToMetricType(mReceiver.MType)

	switch stringType {
	case constants.GaugeType:
		if mReceiver.Value == nil {
			logger.Log.ErrorHTTP(res, errors.New("updating gauge value pointer is nil"), http.StatusBadRequest)
			return
		}
		serv.coreStg.UpdateMetricByName(constants.RenewOperation, stringType, mReceiver.ID, *mReceiver.Value)
		newValue, _ = serv.coreStg.GetMetricByName(constants.GaugeType, mReceiver.ID)
		logger.Log.Info(`Metric "` + mReceiver.ID + `" was successfully updated! New value is ` + strconv.FormatFloat(newValue, 'f', -1, 64))

	case constants.CounterType:
		if mReceiver.Delta == nil {
			logger.Log.ErrorHTTP(res, errors.New("updating counter value pointer is nil"), http.StatusBadRequest)
			return
		}
		serv.coreStg.UpdateMetricByName(constants.AddOperation, stringType, mReceiver.ID, float64(*mReceiver.Delta))
		var counterValue float64
		counterValue, _ = serv.coreStg.GetMetricByName(constants.CounterType, mReceiver.ID)
		newValue = counterValue
		logger.Log.Info(`Metric "` + mReceiver.ID + `" was successfully updated! New value is ` + strconv.Itoa(int(newValue)))

	default:
		convertErr := "ConvertStringToMetricType returns NoneType"
		logger.Log.ErrorHTTP(res, errors.New(convertErr), http.StatusBadRequest)
		return
	}

	respJSON, encodeErr := functionslibrary.EncodeMetricJSON(functionslibrary.ConvertStringToMetricType(mReceiver.MType), mReceiver.ID, newValue)
	if encodeErr != nil {
		logger.Log.ErrorHTTP(res, encodeErr, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(respJSON.Bytes())
}

func (serv *dompserver) WithRequestLog(h http.Handler) http.Handler {
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

func (serv *dompserver) WithResponseLog(h http.Handler) http.Handler {
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
