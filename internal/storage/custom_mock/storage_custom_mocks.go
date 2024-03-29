package storagemockcustom

import "DevOpsMetricsProject/internal/constants"

type StorageMockCustom struct {
	Gauge   map[string]float64
	Counter map[string]int
}

func (ssm *StorageMockCustom) InitMemStorage() {}
func (ssm *StorageMockCustom) ReadMemStorageFields() (g map[string]float64, c map[string]int) {
	return ssm.Gauge, ssm.Counter
}

func (ssm *StorageMockCustom) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) {
}

func (ssm *StorageMockCustom) GetMetricByName(mType constants.MetricType, mName string) (float64, error) {
	return 0, nil
}
