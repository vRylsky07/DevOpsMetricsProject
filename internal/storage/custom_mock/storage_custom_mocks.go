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

func (ssm *StorageMockCustom) UpdateMetricByName(_ constants.UpdateOperation, _ constants.MetricType, _ string, _ float64) {
}

func (ssm *StorageMockCustom) GetMetricByName(_ constants.MetricType, _ string) (float64, error) {
	return 0, nil
}
