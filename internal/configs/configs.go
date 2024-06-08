package configs

import (
	"DevOpsMetricsProject/internal/constants"
	"errors"
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
	RestoreBool   bool   `env:"RESTORE"`
	DatabaseDSN   string `env:"DATABASE_DSN"`
	SaveMode      constants.SaveMode
}

func (cfg *ServerConfig) SetServerConfigFlags() {

	address := flag.String("a", "localhost:8080", "input endpoint address")
	lvl := flag.String("l", "info", "log level")
	interval := flag.Int("i", 300, "metrics save interval")
	temp := flag.String("f", "/tmp/metrics-db.json", "last metrics update")
	restore := flag.Bool("r", true, "restore data or not")
	dsn := flag.String("d", "host=localhost:5432 user=postgres password=123 dbname=testdb sslmode=disable", "database dsn")
	flag.Parse()

	cfg.SaveMode = constants.DatabaseMode
	cfg.Address = *address
	cfg.Loglevel = *lvl
	cfg.StoreInterval = *interval
	cfg.TempFile = *temp
	cfg.RestoreBool = *restore
	cfg.DatabaseDSN = *dsn

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

}

func CreateServerConfig() *ServerConfig {
	cfg := &ServerConfig{}
	cfg.SetServerConfigFlags()

	if cfg.StoreInterval < 0 {
		panic(errors.New("config initialization failed! Store interval value cannot be negative"))
	}

	return cfg
}
