package main

import (
	"DevOpsMetricsProject/internal/server"
	"DevOpsMetricsProject/internal/storage"
	"fmt"
)

func main() {
	storage.InitMemStorage()
	storage.GetMemStrorage().UpdateMetricByName(storage.Counter, "CheckCounter", 15)
	storage.GetMemStrorage().UpdateMetricByName(storage.Gauge, "CheckGauge", 77.6)
	storage.GetMemStrorage().UpdateMetricByName(storage.Counter, "CheckCounter", 14)
	fmt.Println(storage.GetMemStrorage().GetMetricByName(storage.Gauge, "CheckGauge"))
	fmt.Println(storage.GetMemStrorage().GetMetricByName(storage.Counter, "CheckCounter"))
	server.StartServerOnPort(":8080")
}
