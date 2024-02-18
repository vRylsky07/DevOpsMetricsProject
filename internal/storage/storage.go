package storage

import (
	"DevOpsMetricsProject/internal/constants"
	"errors"
	"sort"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage_mocks.go
type StorageInterface interface {
	InitMemStorage()
	ReadMemStorageFields() (g map[string]float64, c map[string]int)
	UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64)
	GetMetricByName(mType constants.MetricType, mName string) (float64, error)
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
func (mStg *MemStorage) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) {

	switch mType {
	case constants.GaugeType:
		if oper == constants.RenewOperation {
			mStg.gauge[mName] = 0
		}
		mStg.gauge[mName] += mValue

	case constants.CounterType:
		if oper == constants.RenewOperation {
			mStg.counter[mName] = 0
		}
		mStg.counter[mName] += int(mValue)
	}
}

// геттер метрик по имени
func (mStg *MemStorage) GetMetricByName(mType constants.MetricType, mName string) (float64, error) {
	switch mType {
	case constants.GaugeType:
		gMetric, wasFound := mStg.gauge[mName]
		if !wasFound {
			return 0.0, errors.New("value was not found")
		}
		return gMetric, nil
	case constants.CounterType:
		cMetric, wasFound := mStg.counter[mName]
		if !wasFound {
			return 0.0, errors.New("value was not found")
		}
		return float64(cMetric), nil
	default:
		return 0.0, errors.New("value was not found")
	}
}

func (mStg *MemStorage) GetSortedKeysArray() (*[]string, *[]string) {
	sortedGauge := []string{}
	sortedCounter := []string{}

	for key := range mStg.gauge {
		sortedGauge = append(sortedGauge, key)
	}

	for key := range mStg.counter {
		sortedCounter = append(sortedCounter, key)
	}

	sort.Slice(sortedGauge, func(i, j int) bool {
		return sortedGauge[i] < sortedGauge[j]
	})

	sort.Slice(sortedCounter, func(i, j int) bool {
		return sortedCounter[i] < sortedCounter[j]
	})

	return &sortedGauge, &sortedCounter
}
