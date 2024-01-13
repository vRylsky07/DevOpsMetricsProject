package main

import server "DevOpsMetricsProject/internal/server"

func main() {
	server.StartServerOnPort(":8080")
}
