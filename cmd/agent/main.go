package main

import (
	"DevOpsMetricsProject/internal/sender"
	"flag"
	"log"
	"sync"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func (cfg *Config) SetConfigFlags() {

	address := flag.String("a", "localhost:8080", "input endpoint address")
	pollInterval := flag.Int("p", 2, "input metrics update interval in seconds")
	reportInterval := flag.Int("r", 10, "input interval to send metrics in seconds")
	flag.Parse()

	cfg.Address = *address
	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateConfig() *Config {
	cfg := &Config{}
	cfg.SetConfigFlags()
	return cfg
}

// запуск http-клиента
func StartAgent() {
	cfg := CreateConfig()

	mSender := sender.CreateSender()
	mSender.SetURL(cfg.Address)

	var wg sync.WaitGroup
	wg.Add(2)
	go mSender.UpdateMetrics(cfg.PollInterval)
	go mSender.SendMetricsHTTP(cfg.ReportInterval)
	wg.Wait()
}

func main() {
	StartAgent()
}
