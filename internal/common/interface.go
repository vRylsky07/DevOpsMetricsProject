package common

import "DevOpsMetricsProject/internal/constants"

type MetricsRepository interface {
	ReadMemStorageFields() (map[string]float64, map[string]int)
	UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) error
	GetMetricByName(mType constants.MetricType, mName string) (float64, error)
	IsValid() bool
}
