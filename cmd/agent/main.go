package main

import (
	//"DevOpsMetricsProject/internal/sender"
	"sync"
)

func main() {
	StartAgent()
}

// запуск http-клиента
func StartAgent() {

	//mSender := sender.CreateSender()

	var wg sync.WaitGroup
	wg.Add(2)
	wg.Wait()
}
