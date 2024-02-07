package main

import (
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func main() {
	StartAgent()
}

// запуск http-клиента
func StartAgent() {
	var wg sync.WaitGroup
	wg.Add(2)
	go sender.UpdateMetrics(2)
	go sender.SendMetricsHTTP(10)
	wg.Wait()
}
