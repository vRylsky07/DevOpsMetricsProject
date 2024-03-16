package agent

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/sender"
	"sync"
)

func Start() {
	cfg := configs.CreateClientConfig()

	mSender := sender.CreateSender()
	mSender.SetDomainURL(cfg.Address)

	var wg sync.WaitGroup
	wg.Add(2)
	go mSender.UpdateMetrics(cfg.PollInterval)
	go mSender.SendMetricsHTTP(cfg.ReportInterval)
	wg.Wait()
}
