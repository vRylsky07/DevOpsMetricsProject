package dompdb

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type DompDB struct {
	db  *sql.DB
	log logger.Recorder
	mtx sync.Mutex
}

func (d *DompDB) PingDB() error {
	return d.db.PingContext(context.TODO())
}

func (d *DompDB) GetAllData() (*map[string]float64, *map[string]int) {
	gaugeOut := make(map[string]float64)
	counterOut := make(map[string]int)

	d.mtx.Lock()
	defer d.mtx.Unlock()

	var rowsG *sql.Rows
	var errG error

	for _, v := range *constants.GetRetryIntervals() {
		if v != 0 {
			d.log.Info("Database is dont resonding (Get gauge data). Retry get data again")
			timer := time.NewTimer(time.Duration(v) * time.Second)
			<-timer.C
		}
		rowsG, errG = d.db.QueryContext(context.TODO(), "SELECT * FROM gauge")
		if errG == nil {
			break
		}
	}

	var rowsC *sql.Rows
	var errC error

	for _, v := range *constants.GetRetryIntervals() {
		if v != 0 {
			d.log.Info("Database is dont resonding (Get counter data). Retry get data again")
			timer := time.NewTimer(time.Duration(v) * time.Second)
			<-timer.C
		}
		rowsC, errC = d.db.QueryContext(context.TODO(), "SELECT * FROM counter")
		if errC == nil {
			break
		}
	}

	if errG != nil || errC != nil {
		return nil, nil
	}

	defer rowsG.Close()

	for rowsG.Next() {
		var name string
		var value float64

		err := rowsG.Scan(&name, &value)

		if err != nil {
			d.log.Error(err.Error())
			return nil, nil
		}

		gaugeOut[name] = value
	}

	err := rowsG.Err()
	if err != nil {
		d.log.Error(err.Error())
		return nil, nil
	}

	defer rowsC.Close()

	for rowsC.Next() {
		var name string
		var value int

		err := rowsC.Scan(&name, &value)

		if err != nil {
			d.log.Error(err.Error())
			return nil, nil
		}

		counterOut[name] = value
	}

	err = rowsC.Err()
	if err != nil {
		d.log.Error(err.Error())
		return nil, nil
	}

	return &gaugeOut, &counterOut
}

func (d *DompDB) IsValid() bool {
	return d.db != nil
}

func NewDompDB(dsn string, log logger.Recorder) (backup.MetricsBackup, error) {
	log.Info("Database DSN: " + dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if errPing := db.PingContext(context.TODO()); errPing != nil {
		return nil, errPing
	}

	errPrep := PrepareTablesDB(db, log)

	if errPrep != nil {
		return nil, errPrep
	}

	dompdb := &DompDB{db: db, log: log}

	log.Info("Connection to Database is OK")
	return dompdb, nil
}

func CheckTableExist(db *sql.DB, log logger.Recorder, tableName string) bool {

	row := db.QueryRowContext(context.TODO(),
		"SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = $1);", tableName)

	isExisted := false

	err := row.Scan(&isExisted)

	if err != nil {
		log.Error(err.Error())
		return false
	}

	return isExisted
}

func PrepareTablesDB(db *sql.DB, log logger.Recorder) error {

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	if !CheckTableExist(db, log, "gauge") {
		_, errCreate := tx.ExecContext(context.TODO(), `CREATE TABLE gauge(
			"name" varchar PRIMARY KEY,
			"value" double precision
			);`)

		if errCreate != nil {
			log.Error(errCreate.Error())
			tx.Rollback()
			return errCreate
		}

		log.Info("Database table named 'gauge' was been successfully created")
	}

	if !CheckTableExist(db, log, "counter") {
		_, errCreate := tx.ExecContext(context.TODO(), `CREATE TABLE counter(
			"name" varchar PRIMARY KEY,
			"value" bigint
			);`)

		if errCreate != nil {
			log.Error(errCreate.Error())
			tx.Rollback()
			return errCreate
		}

		log.Info("Database table named 'counter' was been successfully created")
	}

	return tx.Commit()
}

func (d *DompDB) UpdateMetricBackup(mType constants.MetricType, mName string, mValue float64) error {

	d.mtx.Lock()
	defer d.mtx.Unlock()
	tx, errBegin := d.db.Begin()

	if errBegin != nil {
		return errBegin
	}

	q := fmt.Sprintf(`INSERT INTO %s (name, value)
	VALUES ($1, $2)
	ON CONFLICT (name)
	DO UPDATE SET
	name=EXCLUDED.name,
	value=$3;`, funcslib.ConvertMetricTypeToString(mType))

	_, err := tx.ExecContext(context.TODO(), q, mName, mValue, mValue)

	if err != nil {
		d.log.Error(err.Error())
		tx.Rollback()
		return err
	}

	d.log.Info("Database update metric successfull.", zap.String("MetricName", mName), zap.Float64("Value", mValue))

	return tx.Commit()
}
