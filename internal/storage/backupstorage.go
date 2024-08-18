package storage

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/logger"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type BackupSupportStorage struct {
	log    logger.Recorder
	backup backup.MetricsBackup
	MemStorage
}

func NewBackupSupportStorage(restore bool, backup backup.MetricsBackup, log logger.Recorder) MetricsRepository {
	mStg := &BackupSupportStorage{}
	mStg.backup = backup
	mStg.log = log
	mStg.gauge, mStg.counter = map[string]float64{}, map[string]int{}
	if restore {
		mStg.RestoreDataFromBackup()
	}
	return mStg
}

func (mStg *BackupSupportStorage) IsValid() bool {
	return mStg.backup != nil && mStg.log != nil && mStg.gauge != nil && mStg.counter != nil
}

func (mStg *BackupSupportStorage) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) error {

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

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				continue
			}
		} else {
			break
		}
	}

	return err
}

func (mStg *BackupSupportStorage) RestoreDataFromBackup() {
	g, c := mStg.backup.GetAllData()

	if g == nil || c == nil || (len(*g) == 0) || (len(*c) == 0) {
		mStg.log.Info("Database is empty")
		return
	}

	for k, v := range *g {
		mStg.UpdateMetricByName(constants.RenewOperation, constants.GaugeType, k, v)
	}

	for k, v := range *c {
		mStg.UpdateMetricByName(constants.AddOperation, constants.CounterType, k, float64(v))
	}

	mStg.log.Info("Restore data from database successfully")
}
