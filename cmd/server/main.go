package main

import (
	"DevOpsMetricsProject/internal/server"
	//"DevOpsMetricsProject/internal/storage"
)

func main() {
	server.StartServerOnPort(":8080")
}
