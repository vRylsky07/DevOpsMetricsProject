package storage

import (
	"DevOpsMetricsProject/internal/constants"
	"errors"
	"sync"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage_mocks.go
type StorageInterface interface {
	InitMemStorage()
	ReadMemStorageFields() (g map[string]float64, c map[string]int)
	UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64)
	GetMetricByName(mType constants.MetricType, mName string) (float64, error)
}

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int
	mtx     sync.Mutex
}

func (mStg *MemStorage) ReadMemStorageFields() (g map[string]float64, c map[string]int) {
	gaugeOut := make(map[string]float64)

	mStg.mtx.Lock()
	for k, v := range mStg.gauge {
		gaugeOut[k] = v
	}
	mStg.mtx.Unlock()

	counterOut := make(map[string]int)

	mStg.mtx.Lock()
	for k, v := range mStg.counter {
		counterOut[k] = v
	}
	mStg.mtx.Unlock()

	return gaugeOut, counterOut
}

func (mStg *MemStorage) InitMemStorage() {
	mStg.gauge, mStg.counter = map[string]float64{}, map[string]int{}
}

func (mStg *MemStorage) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) {

	mStg.mtx.Lock()
	defer mStg.mtx.Unlock()
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

func (mStg *MemStorage) GetMetricByName(mType constants.MetricType, mName string) (float64, error) {
	mStg.mtx.Lock()
	defer mStg.mtx.Unlock()

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
