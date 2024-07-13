package server

import (
	"DevOpsMetricsProject/internal/constants"
	funcslib "DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
)

func (serv *dompserver) GetMainPageHandler(res http.ResponseWriter, req *http.Request) {

	if !serv.IsValid() {
		http.Error(res, "Server not working fine please check its initialization!", http.StatusInternalServerError)
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

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(htmlFinal))
}

func (serv *dompserver) GetMetricHandler(res http.ResponseWriter, req *http.Request) {

	if !serv.IsValid() {
		http.Error(res, "Server not working fine please check its initialization!", http.StatusInternalServerError)
		return
	}

	mType := chi.URLParam(req, "mType")
	mName := chi.URLParam(req, "mName")

	if mType != "gauge" && mType != "counter" {
		logger.Log.ErrorHTTP(res, errors.New("your request is incorrect! Metric type should be `gauge` or `counter`"), http.StatusBadRequest)
		return
	}

	mTypeConst := funcslib.ConvertStringToMetricType(mType)

	if mTypeConst == constants.NoneType {
		http.Error(res, "Your request is incorrect! Metric type conversion failed!", http.StatusBadRequest)
		return
	}

	valueToReturn, gettingValueError := serv.coreStg.GetMetricByName(mTypeConst, mName)

	if gettingValueError != nil {
		http.Error(res, "This metric does not exist or was not been updated yet", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(strconv.FormatFloat(valueToReturn, 'f', -1, 64)))
}

func (serv *dompserver) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {

	if !serv.IsValid() {
		http.Error(res, "Server not working fine please check its initialization!", http.StatusInternalServerError)
		return
	}

	mType := chi.URLParam(req, "mType")
	mName := chi.URLParam(req, "mName")
	mValue := chi.URLParam(req, "mValue")

	if mName == "" {
		http.Error(res, "Your request is incorrect! Please enter valid metric name", http.StatusNotFound)
		return
	}

	mTypeConst := funcslib.ConvertStringToMetricType(mType)

	valueInFloat, err := strconv.ParseFloat(mValue, 64)

	if (mType == "gauge" || mType == "counter") && err == nil {

		switch mTypeConst {
		case constants.GaugeType:
			serv.coreStg.UpdateMetricByName(constants.RenewOperation, mTypeConst, mName, valueInFloat)
		case constants.CounterType:
			serv.coreStg.UpdateMetricByName(constants.AddOperation, mTypeConst, mName, valueInFloat)
		}

		switch serv.cfg.SaveMode {
		case constants.DatabaseMode:
			serv.dompdb.UpdateMetricDB(mTypeConst, mName, valueInFloat)
		case constants.FileMode:
			mJSON, errEnc := funcslib.EncodeMetricJSON(mTypeConst, mName, valueInFloat)
			if errEnc == nil {
				serv.SaveCurrentMetrics(mJSON)
			}
		}

		res.Header().Set("Content-Type", "text/html")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Metrics was been updated! Thank you!"))
		return
	}

	http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
}

func (serv *dompserver) IncorrectRequestHandler(res http.ResponseWriter, req *http.Request) {
	if !serv.IsValid() {
		http.Error(res, "Server not working fine please check its initialization!", http.StatusInternalServerError)
		return
	}
	http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
}

func (serv *dompserver) MetricHandlerJSON(res http.ResponseWriter, req *http.Request) {

	if !serv.IsValid() {
		logger.Log.ErrorHTTP(res, errors.New("server not working fine please check its initialization"), http.StatusInternalServerError)
		return
	}

	var body bytes.Buffer
	body.ReadFrom(req.Body)

	isUpdate := (req.URL.Path == "/update" || req.URL.Path == "/update/")

	readCloser := io.NopCloser(strings.NewReader(body.String()))

	mReceiver, err := funcslib.DecodeMetricJSON(readCloser)

	if err != nil {
		logger.Log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	var newValue float64

	mType := funcslib.ConvertStringToMetricType(mReceiver.MType)

	if isUpdate {
		err := funcslib.UpdateStorageInterfaceByMetricStruct(serv.coreStg, mReceiver)
		if err != nil {
			logger.Log.ErrorHTTP(res, err, http.StatusInternalServerError)
		}
	}

	newValue, _ = serv.coreStg.GetMetricByName(mType, mReceiver.ID)

	if isUpdate {
		switch serv.cfg.SaveMode {
		case constants.FileMode:
			updatedJSON, err := funcslib.EncodeMetricJSON(mType, mReceiver.ID, newValue)
			if err == nil {
				serv.SaveCurrentMetrics(updatedJSON)
			}
		case constants.DatabaseMode:
			serv.dompdb.UpdateMetricDB(mType, mReceiver.ID, newValue)
		}
	}

	respJSON, encodeErr := funcslib.EncodeMetricJSON(funcslib.ConvertStringToMetricType(mReceiver.MType), mReceiver.ID, newValue)
	if encodeErr != nil {
		logger.Log.ErrorHTTP(res, encodeErr, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(respJSON.Bytes())
}

func (serv *dompserver) UpdateBatchHandler(res http.ResponseWriter, req *http.Request) {
	if !serv.IsValid() {
		logger.Log.ErrorHTTP(res, errors.New("server not working fine please check its initialization"), http.StatusInternalServerError)
		return
	}

	var body bytes.Buffer
	body.ReadFrom(req.Body)

	readCloser := io.NopCloser(strings.NewReader(body.String()))

	mReceiver, err := funcslib.DecodeBatchJSON(readCloser)

	if err != nil {
		logger.Log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	for _, m := range *mReceiver {
		err = funcslib.UpdateStorageInterfaceByMetricStruct(serv.coreStg, &m)
		if err != nil {
			logger.Log.ErrorHTTP(res, err, http.StatusInternalServerError)
			return
		}
	}

	err = serv.dompdb.UpdateBatchDB(serv.coreStg)

	if err != nil {
		logger.Log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Metrics was been updated by batch! Thank you!"))
}

func (serv *dompserver) PingDatabaseHandler(res http.ResponseWriter, req *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := serv.dompdb.db.PingContext(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
}

// Middlewares

func (serv *dompserver) WithRequestLog(h http.Handler) http.Handler {
	if !serv.IsValid() {
		return h
	}
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
	if !serv.IsValid() {
		return h
	}

	logFn := func(w http.ResponseWriter, r *http.Request) {

		rlw := &ResponceLogWriter{w, 0, 0}
		h.ServeHTTP(rlw, r)

		logger.Log.Info("Server responding", zap.Int("status", rlw.statusCode), zap.Int("size", rlw.size))
	}

	return http.HandlerFunc(logFn)
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			logger.Log.ErrorHTTP(w, err, http.StatusNotFound)
			return
		}

		allowTypes := &[]string{"application/json", "text/html"}

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz, AllowedTypes: allowTypes}, r)
	})
}

func DecompressHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			logger.Log.Info("No decompression")
			next.ServeHTTP(w, r)
			return
		}

		data, err := funcslib.DecompressData(r.Body)

		if err != nil {
			logger.Log.Info("DecompressData failed!")
			next.ServeHTTP(w, r)
			return
		}

		r.Body = data

		logger.Log.Info("Using GZIP decompression")
		next.ServeHTTP(w, r)
	})
}
