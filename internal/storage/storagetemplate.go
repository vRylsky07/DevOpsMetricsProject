package storage

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
)

//go:generate mockgen -source=storagetemplate.go -destination=mocks/storage_mocks.go
type MetricsRepository interface {
	InitMemStorage(restore bool, backup backup.MetricsBackup, log logger.Recorder)
	ReadMemStorageFields() (g map[string]float64, c map[string]int)
	UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64)
	GetMetricByName(mType constants.MetricType, mName string) (float64, error)
	CheckBackupStatus() error
}
