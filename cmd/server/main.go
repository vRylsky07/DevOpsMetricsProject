package main

import (
	"DevOpsMetricsProject/internal/server"
	"flag"
)

func main() {
	address := flag.String("a", "localhost:8080", "input endpoint address")
	flag.Parse()

	server.StartServerOnPort(*address)
}
