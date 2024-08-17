package storage

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
	"errors"
	"sync"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int
	mtx     sync.Mutex
}

func (mStg *MemStorage) CheckBackupStatus() error {
	return nil
}

func (mStg *MemStorage) ReadMemStorageFields() (g map[string]float64, c map[string]int) {
	gaugeOut := make(map[string]float64)
	counterOut := make(map[string]int)

	mStg.mtx.Lock()
	defer mStg.mtx.Unlock()
	for k, v := range mStg.gauge {
		gaugeOut[k] = v
	}

	for k, v := range mStg.counter {
		counterOut[k] = v
	}

	return gaugeOut, counterOut
}

func (mStg *MemStorage) InitMemStorage(_ bool, _ backup.MetricsBackup, _ logger.Recorder) {
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
