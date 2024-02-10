package main

import (
	"DevOpsMetricsProject/internal/server"
)

func main() {
	server.StartServerOnPort(":8080")
}
