package main

import (
	"DevOpsMetricsProject/internal/metrics"
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func main() {
	StartAgent()
}

// запуск http-клиента
func StartAgent() {

	mCollector := &metrics.MetricsCollector{}
	mCollector.InitMetricsCollector()

	mSender := sender.CreateSender()

	var wg sync.WaitGroup
	wg.Add(2)
	mSender.UpdateMetrics(mCollector, 2)
	wg.Wait()
}
