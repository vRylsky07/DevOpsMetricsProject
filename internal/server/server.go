package server

import (
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/storage"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type dompserver struct {
	coreMux *chi.Mux
	coreStg *storage.MemStorage
}

// главная функция запуска и инициализации сервера
func StartServerOnPort(port string) {
	dompserv := CreateNewServer()
	if dompserv.coreMux == nil || dompserv.coreStg == nil {
		return
	}
	err := http.ListenAndServe(port, dompserv.coreMux)
	if err != nil {
		panic(err)
	}
}

// создание и настройка нового маршрутизатора
func CreateNewServer() *dompserver {
	coreMux := chi.NewRouter()
	coreStg := &storage.MemStorage{}
	coreStg.InitMemStorage()
	dompserv := &dompserver{coreMux: coreMux, coreStg: coreStg}

	coreMux.HandleFunc("/update/", dompserv.UpdateMetricHandler)
	coreMux.Route("/update", func(r chi.Router) {
		r.Post("/{mType}/{mName}/{mValue}", dompserv.UpdateMetricHandler)
		r.Get("/{mType}/{mName}/{mValue}", func(res http.ResponseWriter, req *http.Request) {
			http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
		})
	})
	coreMux.Route("/value", func(r chi.Router) {
		r.Get("/{mName}/{mValue}", dompserv.GetMetricHandler)
	})
	return dompserv
}

func (serv *dompserver) GetMetricHandler(res http.ResponseWriter, req *http.Request) {

	mType := chi.URLParam(req, "mType")
	mName := chi.URLParam(req, "mName")

	if (mType != "gauge" && mType != "counter") || serv == nil || serv.coreMux == nil || serv.coreStg == nil {
		http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
		return
	}

	mTypeConst := functionslibrary.ConvertStringToMetricType(mType)

	if mTypeConst == -1 {
		http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
		return
	}

	valueToReturn, gettingValueError := serv.coreStg.GetMetricByName(mTypeConst, mName)

	if gettingValueError == nil {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(strconv.FormatFloat(valueToReturn, 'f', 2, 64)))
		return
	} else {
		http.Error(res, "Your request is incorrect! Please enter valid metric name", http.StatusNotFound)
		return
	}
}

// хэндлер POST-запроса на /update/
func (serv *dompserver) UpdateMetricHandler(res http.ResponseWriter, req *http.Request) {

	mType := chi.URLParam(req, "mType")
	mName := chi.URLParam(req, "mName")
	mValue := chi.URLParam(req, "mValue")

	if mName == "" {
		http.Error(res, "Your request is incorrect! Please enter valid metric name", http.StatusNotFound)
		return
	}

	_, err := strconv.ParseFloat(mValue, 64)

	if (mType == "gauge" || mType == "counter") && err == nil {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte("Metrics was been updated! Thank you!"))
		return
	}

	http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
}
