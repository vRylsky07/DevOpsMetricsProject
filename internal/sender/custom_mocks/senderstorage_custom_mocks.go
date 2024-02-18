package senderstoragecustommock

import "DevOpsMetricsProject/internal/constants"

type SenderStorageMock struct {
	Gauge   map[string]float64
	Counter map[string]int
}

func (ssm *SenderStorageMock) InitMemStorage() {}
func (ssm *SenderStorageMock) ReadMemStorageFields() (g map[string]float64, c map[string]int) {
	return ssm.Gauge, ssm.Counter
}

func (ssm *SenderStorageMock) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) {
}

func (ssm *SenderStorageMock) GetMetricByName(mType constants.MetricType, mName string) (float64, error) {
	return 0, nil
}
