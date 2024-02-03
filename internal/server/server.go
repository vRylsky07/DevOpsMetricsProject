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
var dompMux *http.ServeMux

// инициализация маршрутизатора
func CoreMuxInit() {
	dompMux = http.NewServeMux()
	dompMux.HandleFunc("/update/", mainPage)
}

// геттер маршрутизатора
func GetDompMux() *http.ServeMux {
	return dompMux
}

// хэндлер POST-запроса на /update/
func mainPage(res http.ResponseWriter, req *http.Request) {
	finalReqStatus := http.StatusBadRequest

	if req.Method == http.MethodPost {
		/*if isValid := HeadersValidator(&req.Header); !isValid {
			http.Error(res, "Headers validation FAILED", finalReqStatus)
			return
		}
		*/
		parsedURL := ParseAndValidateURL(string(req.URL.Path), &finalReqStatus)
		if parsedURL != nil {
			res.WriteHeader(finalReqStatus)
			res.Write([]byte("Metrics was been updated! Thank you!"))
			return
		}
	}
	http.Error(res, "Your request is incorrect!", finalReqStatus)
}

// Разбивка URL-строки и проверка соответствия передаваемых данных
func ParseAndValidateURL(p string, status *int) []string {
	var pathParsed = strings.Split(string(p), "/")

	if len(pathParsed) <= 3 || (len(pathParsed) >= 4 && pathParsed[3] == "") {
		*status = http.StatusNotFound
		return nil
	}

	if len(pathParsed) == 5 && pathParsed[1] == "update" && (pathParsed[2] == "gauge" || pathParsed[2] == "counter") {

		_, valueValidate := strconv.ParseFloat(pathParsed[4], 64)

		if valueValidate == nil {
			*status = http.StatusOK
			return pathParsed
		}
	}
	return nil
}

// валидатор хедера Content-type
/*func HeadersValidator(headers *http.Header) bool {
	headerContentType, ok := (*headers)["Content-Type"]

	if ok && FindHeaderValue(headerContentType, "text/plain") {
		return true
	}
	return false
}


func FindHeaderValue(arrayStr []string, comparingStr string) bool {
	for _, value := range arrayStr {
		if value == comparingStr {
			return true
		}
	}
	return false
}
*/
