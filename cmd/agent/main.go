package main

import (
	"DevOpsMetricsProject/internal/sender"
	"flag"
	"sync"
)

func main() {
	address := flag.String("a", "localhost:8080", "input endpoint address")
	pollInterval := flag.Int("p", 2, "input metrics update interval in seconds")
	reportInterval := flag.Int("r", 10, "input interval to send metrics in seconds")
	flag.Parse()

	StartAgent(*address, *pollInterval, *reportInterval)
}

// запуск http-клиента
func StartAgent(address string, pollInterval int, reportInterval int) {
	mSender := sender.CreateSender()
	mSender.SetURL(address)

	var wg sync.WaitGroup
	wg.Add(2)
	go mSender.UpdateMetrics(pollInterval)
	go mSender.SendMetricsHTTP(reportInterval)
	wg.Wait()
}
