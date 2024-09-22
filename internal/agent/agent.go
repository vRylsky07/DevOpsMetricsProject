package agent

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/coretypes"
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func Start() {
	cfg := configs.CreateClientConfig()
	jobs := make(chan *coretypes.ReqProps, 10)

	mSender, err := sender.CreateSender(cfg, jobs)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)
	for i := 1; i <= 1; i++ {
		go mSender.RequestSendingWorker(i, jobs)
	}
	go mSender.UpdateMetrics()
	go mSender.SendMetricsHTTP()
	mSender.GetLogger().Info("Agent was successfully started!")
	wg.Wait()
}
