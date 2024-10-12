package agent

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func Start() {
	cfg := configs.CreateClientConfig()

	mSender, err := sender.CreateSender(cfg)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go mSender.UpdateMetrics()
	go mSender.SendMetricsHTTP()
	mSender.GetLogger().Info("Agent was successfully started!")
	wg.Wait()
}
