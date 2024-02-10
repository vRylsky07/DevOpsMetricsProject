package main

import (
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func main() {
	StartAgent(2, 10)
}

// запуск http-клиента
func StartAgent(pollInterval int, reportInterval int) {

	mSender := sender.CreateSender()

	var wg sync.WaitGroup
	wg.Add(2)
	go mSender.UpdateMetrics(pollInterval)
	go mSender.SendMetricsHTTP(reportInterval)
	wg.Wait()
}
