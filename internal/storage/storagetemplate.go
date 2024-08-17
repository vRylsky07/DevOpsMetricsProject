package storage

import (
	"DevOpsMetricsProject/internal/constants"
)

//go:generate mockgen -source=storagetemplate.go -destination=mocks/storage_mocks.go
type MetricsRepository interface {
	ReadMemStorageFields() (map[string]float64, map[string]int)
	UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64)
	GetMetricByName(mType constants.MetricType, mName string) (float64, error)
	IsValid() bool
	CheckBackupStatus() error
}
