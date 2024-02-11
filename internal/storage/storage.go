package storage

import (
	"DevOpsMetricsProject/internal/constants"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage_mocks.go
type StorageInterface interface {
	InitMemStorage()
	ReadMemStorageFields() (g map[string]float64, c map[string]int)
	UpdateMetricByName(mType constants.MetricType, mName string, mValue float64)
	GetMetricByName(mType constants.MetricType, mName string) float64
}

// описание хранилища данных
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int
}

func (mStg *MemStorage) ReadMemStorageFields() (g map[string]float64, c map[string]int) {
	return mStg.gauge, mStg.counter
}

func (mStg *MemStorage) InitMemStorage() {
	mStg.gauge, mStg.counter = map[string]float64{}, map[string]int{}
}

// обновление метрик с указанием Enum типа
func (mStg *MemStorage) UpdateMetricByName(mType constants.MetricType, mName string, mValue float64) {

	switch mType {
	case constants.GaugeType:
		mStg.gauge[mName] = mValue
	case constants.CounterType:
		mStg.counter[mName] += int(mValue)
	}
}

// геттер метрик по имени
func (mStg *MemStorage) GetMetricByName(mType constants.MetricType, mName string) float64 {
	switch mType {
	case constants.GaugeType:
		gMetric, wasFound := mStg.gauge[mName]
		if !wasFound {
			return 0.0
		}
		return gMetric
	case constants.CounterType:
		cMetric, wasFound := mStg.counter[mName]
		if !wasFound {
			return 0.0
		}
		return float64(cMetric)
	default:
		return 0.0
	}
}
