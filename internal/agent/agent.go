package agent

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func Start() {
	cfg := configs.CreateClientConfig()
	logger.Initialize(cfg.Loglevel, "agent_")

	mSender := sender.CreateSender(cfg)

	var wg sync.WaitGroup
	wg.Add(2)
	go mSender.UpdateMetrics()
	go mSender.SendMetricsHTTP()
	logger.Log.Info("Agent was successfully started!")
	wg.Wait()
}
