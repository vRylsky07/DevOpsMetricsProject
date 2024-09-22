package backup

import "DevOpsMetricsProject/internal/constants"

//go:generate mockgen -source=backup.go -destination=mocks/backup_mocks.go
type MetricsBackup interface {
	UpdateMetricBackup(mType constants.MetricType, mName string, mValue float64) error
	IsValid() bool
	GetAllData() (*map[string]float64, *map[string]int)
}

type PingerDB interface {
	PingDB() error
}
