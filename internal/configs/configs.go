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
}

func (cfg *ClientConfig) SetClientConfigFlags() {

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

func CreateClientConfig() *ClientConfig {
	cfg := &ClientConfig{}
	cfg.SetClientConfigFlags()
	return cfg
}

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

// SERVER CONFIG

func (cfg *ServerConfig) SetServerConfigFlags() {

	address := flag.String("a", "localhost:8080", "input endpoint address")
	flag.Parse()

	cfg.Address = *address

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
