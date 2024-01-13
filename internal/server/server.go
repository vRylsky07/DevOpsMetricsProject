package server

import "net/http"

func StartServerOnPort(port string) {
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
