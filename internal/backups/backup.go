package backup

import "DevOpsMetricsProject/internal/constants"

type MetricsBackup interface {
	UpdateMetricDB(mType constants.MetricType, mName string, mValue float64) error
	UpdateBatchDB(gauge *map[string]float64, counter *map[string]int) error
	IsValid() bool
	GetAllData() (g map[string]float64, c map[string]int)
	CheckBackupStatus() error
}
