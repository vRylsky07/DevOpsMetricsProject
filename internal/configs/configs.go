package configs

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v10"
)

// CLIENT CONFIG

type ClientConfig struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	Loglevel       string `env:"LOG_LEVEL"`
}

func (cfg *ClientConfig) SetClientConfigFlags() {

	address := flag.String("a", "localhost:8080", "input endpoint address")
	pollInterval := flag.Int("p", 2, "input metrics update interval in seconds")
	reportInterval := flag.Int("r", 10, "input interval to send metrics in seconds")
	lvl := flag.String("l", "info", "log level")
	flag.Parse()

	cfg.Address = *address
	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.Loglevel = *lvl

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateClientConfig() *ClientConfig {
	cfg := &ClientConfig{}
	cfg.SetClientConfigFlags()
	return cfg
}

// SERVER CONFIG

type ServerConfig struct {
	Address       string `env:"ADDRESS"`
	Loglevel      string `env:"LOG_LEVEL"`
	StoreInterval int    `env:"STORE_INTERVAL"`
	TempFile      string `env:"FILE_STORAGE_PATH"`
}

func (cfg *ServerConfig) SetServerConfigFlags() {

	address := flag.String("a", "localhost:8080", "input endpoint address")
	lvl := flag.String("l", "info", "log level")
	interval := flag.Int("i", 300, "metrics save interval")
	temp := flag.String("f", "metrics-db.json", "last metrics update")
	flag.Parse()

	cfg.Address = *address
	cfg.Loglevel = *lvl
	cfg.StoreInterval = *interval
	cfg.TempFile = *temp

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.SetServerConfigFlags()
	return cfg
}
