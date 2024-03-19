package agent

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func Start() {
	cfg := configs.CreateClientConfig()
	logger.Initialize(cfg.Loglevel)

	mSender := sender.CreateSender()
	mSender.SetDomainURL(cfg.Address)

	var wg sync.WaitGroup
	wg.Add(2)
	go mSender.UpdateMetrics(cfg.PollInterval)
	go mSender.SendMetricsHTTP(cfg.ReportInterval)
	logger.Log.Info("Agent was successfully started!")
	wg.Wait()
}
