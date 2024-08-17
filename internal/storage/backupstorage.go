package storage

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
	"errors"
	"sync"
	"time"
)

type BackupSupportStorage struct {
	backup  backup.MetricsBackup
	log     logger.Recorder
	gauge   map[string]float64
	counter map[string]int
	mtx     sync.Mutex
}

func (mStg *BackupSupportStorage) CheckBackupStatus() error {
	return mStg.backup.CheckBackupStatus()
}

func (mStg *BackupSupportStorage) ReadMemStorageFields() (g map[string]float64, c map[string]int) {
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

func (mStg *BackupSupportStorage) InitMemStorage(restore bool, backup backup.MetricsBackup, log logger.Recorder) {
	mStg.backup = backup
	mStg.log = log
	mStg.gauge, mStg.counter = map[string]float64{}, map[string]int{}
	if restore {
		mStg.RestoreDataFromBackup()
	}
}

func (mStg *BackupSupportStorage) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) {

	mStg.mtx.Lock()
	defer mStg.mtx.Unlock()
	var updatedValue float64

	switch mType {
	case constants.GaugeType:
		if oper == constants.RenewOperation {
			mStg.gauge[mName] = 0
		}
		mStg.gauge[mName] += mValue
		updatedValue = mStg.gauge[mName]

	case constants.CounterType:
		if oper == constants.RenewOperation {
			mStg.counter[mName] = 0
		}
		mStg.counter[mName] += int(mValue)
		updatedValue = float64(mStg.counter[mName])
	}

	var err error
	for _, v := range *constants.GetRetryIntervals() {
		if v != 0 {
			mStg.log.Info("Update database or current metrics file failed. Try again...")
			err = nil
			timer := time.NewTimer(time.Duration(v) * time.Second)
			<-timer.C
		}
		err = mStg.backup.UpdateMetricDB(mType, mName, updatedValue)
		if err == nil {
			mStg.log.Info("Backup storage was successfully updated")
			break
		}
	}
}

func (mStg *BackupSupportStorage) GetMetricByName(mType constants.MetricType, mName string) (float64, error) {
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

func (mStg *BackupSupportStorage) RestoreDataFromBackup() {
	g, c := mStg.backup.GetAllData()

	if g == nil || c == nil || (len(g) == 0) || (len(c) == 0) {
		mStg.log.Info("Database is empty")
		return
	}

	for k, v := range g {
		mStg.UpdateMetricByName(constants.RenewOperation, constants.GaugeType, k, v)
	}

	for k, v := range c {
		mStg.UpdateMetricByName(constants.AddOperation, constants.CounterType, k, float64(v))
	}

	mStg.log.Info("Restore data from database successfully")
}
