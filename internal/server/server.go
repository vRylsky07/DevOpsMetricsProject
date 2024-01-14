package server

import (
	"net/http"
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
	if req.Method == http.MethodPost {
		parsedURL := ParsePath(string(req.URL.Path))
		if parsedURL != nil {
			res.Write([]byte("Do: " + parsedURL[1]))
			res.Write([]byte(" Type: " + parsedURL[2]))
			return
		}
		res.Write([]byte("Запрос составлен неверно. Проверьте и отправьте еще раз."))
	} else {
		res.Write(testFormPOST)
	}
}

// Разбивка строки и проверка соответствия передаваемых данных
func ParsePath(p string) []string {
	var pathParsed = strings.Split(string(p), "/")
	if len(pathParsed) == 5 && pathParsed[1] == "update" && (pathParsed[2] == "gauge" || pathParsed[2] == "counter") {
		return pathParsed
	}
	return nil
}

// тестовая страница для отправки POST запроса
var testFormPOST []byte = []byte(`<html>
<head>
<title></title>
</head>
<body>
	<h3> Страница тестирования POST-запросов </h3> 
	<h3>Введите логин и пароль для отправки: </h3>
	<form action="/update/counter/Metric/36.6" method="post">
		<label>Логин <input type="text" name="login"></label>
		<label>Пароль <input type="password" name="password"></label>
		<input type="submit" value="Login">
	</form>
</body>
</html>`)
