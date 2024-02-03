package main

import (
	"DevOpsMetricsProject/internal/server"
	"DevOpsMetricsProject/internal/storage"
)

var A11 int = 15

func main() {
	storage.InitMemStorage()
	server.StartServerOnPort(":8080")
}
