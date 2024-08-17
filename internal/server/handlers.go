package server

import (
	"DevOpsMetricsProject/internal/constants"
	funcslib "DevOpsMetricsProject/internal/funcslib"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		serv.log.ErrorHTTP(res, errors.New("your request is incorrect! Metric type should be `gauge` or `counter`"), http.StatusBadRequest)
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

	RetryLoopLabel:
		for _, v := range *constants.GetRetryIntervals() {
			if v != 0 {
				serv.log.Info("Update database or current metrics file failed. Try again...")
				timer := time.NewTimer(time.Duration(v) * time.Second)
				<-timer.C
			}
			var err error
			switch serv.cfg.SaveMode {
			case constants.DatabaseMode:
				err = serv.db.UpdateMetricDB(mTypeConst, mName, valueInFloat)
				if err == nil {
					break RetryLoopLabel
				}
			case constants.FileMode:
				var mJSON *bytes.Buffer
				var errJSON error
				mJSON, errJSON = funcslib.EncodeMetricJSON(mTypeConst, mName, valueInFloat)

				if errJSON != nil {
					serv.log.Error(errJSON.Error())
					break RetryLoopLabel
				}

				err = serv.SaveCurrentMetrics(mJSON)
				if err == nil {
					if v != 0 {
						serv.log.Info("Database or current metrics was successfully updated")
					}
					break RetryLoopLabel
				}
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
		serv.log.ErrorHTTP(res, errors.New("server not working fine please check its initialization"), http.StatusInternalServerError)
		return
	}

	var body bytes.Buffer
	body.ReadFrom(req.Body)

	isUpdate := (req.URL.Path == "/update" || req.URL.Path == "/update/")

	readCloser := io.NopCloser(strings.NewReader(body.String()))

	mReceiver, err := funcslib.DecodeMetricJSON(readCloser)

	if err != nil {
		serv.log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	var newValue float64

	mType := funcslib.ConvertStringToMetricType(mReceiver.MType)

	if isUpdate {
		err := serv.coreStg.UpdateMetricByName()
		if err != nil {
			serv.log.ErrorHTTP(res, err, http.StatusInternalServerError)
		}
	}

	newValue, _ = serv.coreStg.GetMetricByName(mType, mReceiver.ID)

	if isUpdate {
	RetryLoopLabel:
		for _, v := range *constants.GetRetryIntervals() {
			if v != 0 {
				serv.log.Info("Update database or current metrics file failed. Try again...")
				timer := time.NewTimer(time.Duration(v) * time.Second)
				<-timer.C
			}
			var err error
			switch serv.cfg.SaveMode {
			case constants.FileMode:
				updatedJSON, errJSON := funcslib.EncodeMetricJSON(mType, mReceiver.ID, newValue)
				if errJSON == nil {
					err = serv.SaveCurrentMetrics(updatedJSON)
				}
			case constants.DatabaseMode:
				err = serv.db.UpdateMetricDB(mType, mReceiver.ID, newValue)
			}
			if err == nil {
				if v != 0 {
					serv.log.Info("Database or current metrics was successfully updated")
				}
				break RetryLoopLabel
			}
		}
	}

	respJSON, encodeErr := funcslib.EncodeMetricJSON(funcslib.ConvertStringToMetricType(mReceiver.MType), mReceiver.ID, newValue)
	if encodeErr != nil {
		serv.log.ErrorHTTP(res, encodeErr, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write(respJSON.Bytes())
}

func (serv *dompserver) UpdateBatchHandler(res http.ResponseWriter, req *http.Request) {
	if !serv.IsValid() {
		serv.log.ErrorHTTP(res, errors.New("server not working fine please check its initialization"), http.StatusInternalServerError)
		return
	}

	var body bytes.Buffer
	body.ReadFrom(req.Body)

	readCloser := io.NopCloser(strings.NewReader(body.String()))

	mReceiver, err := funcslib.DecodeBatchJSON(readCloser)

	if err != nil {
		serv.log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	for _, m := range *mReceiver {
		err = funcslib.UpdateStorageInterfaceByMetricStruct(serv.coreStg, &m)
		if err != nil {
			serv.log.ErrorHTTP(res, err, http.StatusInternalServerError)
			return
		}
	}

	for _, v := range *constants.GetRetryIntervals() {
		if v != 0 {
			serv.log.Info("Update database with batches failed. Try again...")
			timer := time.NewTimer(time.Duration(v) * time.Second)
			<-timer.C
		}

		err := serv.db.UpdateBatchDB(serv.coreStg)

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				continue
			}
		}
		break
	}

	if err != nil {
		serv.log.ErrorHTTP(res, err, http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte("Metrics was been updated by batch! Thank you!"))
}

func (serv *dompserver) PingDatabaseHandler(res http.ResponseWriter, _ *http.Request) {

	if err := serv.coreStg.; err != nil {
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

		serv.log.Info("Server got HTTP-request", zap.String("uri", uri), zap.String("method", method), zap.Duration("time", duration))

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

		serv.log.Info("Server responding", zap.Int("status", rlw.statusCode), zap.Int("size", rlw.size))
	}

	return http.HandlerFunc(logFn)
}

func (serv *dompserver) gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			serv.log.ErrorHTTP(w, err, http.StatusNotFound)
			return
		}

		allowTypes := &[]string{"application/json", "text/html"}

		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz, AllowedTypes: allowTypes}, r)
	})
}

func (serv *dompserver) DecompressHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			serv.log.Info("No decompression")
			next.ServeHTTP(w, r)
			return
		}

		data, err := funcslib.DecompressData(r.Body)

		if err != nil {
			serv.log.Info("DecompressData failed!")
			next.ServeHTTP(w, r)
			return
		}

		r.Body = data

		serv.log.Info("Using GZIP decompression")
		next.ServeHTTP(w, r)
	})
}
