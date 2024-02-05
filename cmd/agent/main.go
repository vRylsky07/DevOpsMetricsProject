package main

import (
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go sender.UpdateMetrics(2)
	wg.Wait()
}
