package storagemockcustom

import "DevOpsMetricsProject/internal/constants"

type StorageMockCustom struct {
	Gauge   map[string]float64
	Counter map[string]int
}

func (ssm *StorageMockCustom) ReadMemStorageFields() (map[string]float64, map[string]int) {
	return ssm.Gauge, ssm.Counter
}

func (ssm *StorageMockCustom) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) error {
	return nil
}

func (ssm *StorageMockCustom) GetMetricByName(_ constants.MetricType, _ string) (float64, error) {
	return 0, nil
}

func (ssm *StorageMockCustom) IsValid() bool {
	return false
}

func (ssm *StorageMockCustom) CheckBackupStatus() error {
	return nil
}
