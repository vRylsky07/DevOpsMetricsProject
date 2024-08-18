package backup

import "DevOpsMetricsProject/internal/constants"

type MetricsBackup interface {
	UpdateMetricDB(mType constants.MetricType, mName string, mValue float64) error
	IsValid() bool
	GetAllData() (*map[string]float64, *map[string]int)
}

type PingerDB interface {
	PingDB() error
}
