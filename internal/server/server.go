package server

import (
	"net/http"
	"strconv"
	"strings"
)

// главная функция запуска и инициализации сервера
func StartServerOnPort(port string) {
	CoreMuxInit()
	err := http.ListenAndServe(port, GetDompMux())
	if err != nil {
		panic(err)
	}
}

// инстанс главного в проекте маршрутизатора
var domp_mux *http.ServeMux

// инициализация маршрутизатора
func CoreMuxInit() {
	domp_mux = http.NewServeMux()
	domp_mux.HandleFunc("/update/", mainPage)
}

// геттер маршрутизатора
func GetDompMux() *http.ServeMux {
	return domp_mux
}

// хэндлер POST-запроса на /update/
func mainPage(res http.ResponseWriter, req *http.Request) {
	finalReqStatus := http.StatusBadRequest
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")

	if req.Method == http.MethodPost {
		if isValid, error := HeadersValidator(&req.Header); !isValid {
			http.Error(res, "Headers validation FAILED\n"+error, finalReqStatus)
			return
		}
		parsedURL := ParsePath(string(req.URL.Path), &finalReqStatus)
		if parsedURL != nil {
			res.WriteHeader(finalReqStatus)
			res.Write([]byte("Metrics was been updated! Thank you!"))
			return
		}
	}
	http.Error(res, "Your request is incorrect!", finalReqStatus)
}

// Разбивка строки и проверка соответствия передаваемых данных
func ParsePath(p string, status *int) []string {
	var pathParsed = strings.Split(string(p), "/")

	if len(pathParsed) <= 3 || (len(pathParsed) >= 4 && pathParsed[3] == "") {
		*status = http.StatusNotFound
		return nil
	}

	if len(pathParsed) == 5 && pathParsed[1] == "update" && (pathParsed[2] == "gauge" || pathParsed[2] == "counter") {
		*status = http.StatusOK
		return pathParsed
	}
	return nil
}

// валидатор хедера Content-type
func HeadersValidator(headers *http.Header) (bool, string) {
	headerContentType, ok := (*headers)["Content-Type"]
	if !ok {
		var allHeaders string = "All headers in your request: "
		for header := range *headers {
			allHeaders += header + " "
		}
		return false, "Content-type header not exist " + "\n" + allHeaders
	}

	isValid := len(headerContentType) == 1 && headerContentType[0] == "text/plain"
	var err string

	if !isValid {
		notEmpty := len(headerContentType) == 1
		equalType := headerContentType[0] == "text/plain; charset=utf-8"
		err = "Header content-type not empty? " + strconv.FormatBool(notEmpty) + " Is valid? " + strconv.FormatBool(equalType)
	}

	return isValid, err
}
