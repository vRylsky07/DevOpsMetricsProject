package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/storage"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi"
	"golang.org/x/exp/maps"
)

type dompserver struct {
	coreMux *chi.Mux
	coreStg storage.StorageInterface
}

func Start() {
	cfg := configs.CreateServerConfig()

	dompserv := CreateNewServer()
	if dompserv.coreMux == nil || dompserv.coreStg == nil {
		return
	}
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

	coreMux.Get("/", dompserv.GetMainPageHandler)
	coreMux.Route("/update", func(r chi.Router) {
		r.Post("/{mType}/{mName}/{mValue}", dompserv.UpdateMetricHandler)
		r.Get("/{mType}/{mName}/{mValue}", func(res http.ResponseWriter, req *http.Request) {
			http.Error(res, "Your request is incorrect!", http.StatusBadRequest)
		})
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

	io.WriteString(res, htmlFinal)

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
	}

	if serv == nil {
		http.Error(res, "ERROR! Server not working fine please check its initialization! Serv.coreMux is nil", http.StatusBadRequest)
	}

	if serv == nil {
		http.Error(res, "ERROR! Server not working fine please check its initialization! Serv.coreStg is nil", http.StatusBadRequest)
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
