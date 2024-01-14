package main

import (
	"DevOpsMetricsProject/internal/server"
	"DevOpsMetricsProject/internal/storage"
)

func main() {
	storage.InitMemStorage()
	server.StartServerOnPort(":8080")
}
