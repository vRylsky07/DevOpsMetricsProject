package server

import (
	"DevOpsMetricsProject/internal/constants"
	funcslib "DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"context"
	"database/sql"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type DompInterfaceDB interface {
	UpdateMetricDB(db *sql.DB, mType constants.MetricType, mName string, mValue float64) error
	UpdateBatchDB(db *sql.DB, sStg storage.StorageInterface) error
	IsValid() bool
	GetAllData() (g map[string]float64, c map[string]int)
}

type dompdb struct {
	db  *sql.DB
	mtx sync.Mutex
}

func (d *dompdb) GetAllData() (g map[string]float64, c map[string]int) {
	gaugeOut := make(map[string]float64)
	counterOut := make(map[string]int)

	d.mtx.Lock()
	defer d.mtx.Unlock()

	rowsG, errG := d.db.QueryContext(context.TODO(), "SELECT * FROM gauge")
	rowsC, errC := d.db.QueryContext(context.TODO(), "SELECT * FROM counter")

	if errG != nil || errC != nil {
		return nil, nil
	}

	defer rowsG.Close()

	for rowsG.Next() {
		var name string
		var value float64

		err := rowsG.Scan(&name, &value)

		if err != nil {
			logger.Log.Error(err.Error())
			return nil, nil
		}

		gaugeOut[name] = value
	}

	err := rowsG.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil
	}

	defer rowsC.Close()

	for rowsC.Next() {
		var name string
		var value int

		err := rowsC.Scan(&name, &value)

		if err != nil {
			logger.Log.Error(err.Error())
			return nil, nil
		}

		counterOut[name] = value
	}

	err = rowsC.Err()
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, nil
	}

	return gaugeOut, counterOut
}

func (d *dompdb) IsValid() bool {
	if d == nil || d.db == nil {
		return false
	}
	return true
}

func RunDB(dsn string) (*dompdb, error) {
	logger.Log.Info("Database DSN: " + dsn)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if errPing := db.PingContext(context.TODO()); errPing != nil {
		return nil, errPing
	}

	errPrep := PrepareTablesDB(db)

	if errPrep != nil {
		return nil, errPrep
	}

	dompdb := &dompdb{db: db}

	logger.Log.Info("Connection to Database is OK")
	return dompdb, nil
}

func CheckTableExist(db *sql.DB, tableName string) bool {

	row := db.QueryRowContext(context.TODO(),
		"SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = $1);", tableName)

	isExisted := false

	err := row.Scan(&isExisted)

	if err != nil {
		logger.Log.Error(err.Error())
		return false
	}

	return isExisted
}

func PrepareTablesDB(db *sql.DB) error {

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	if !CheckTableExist(db, "gauge") {
		_, errCreate := tx.ExecContext(context.TODO(), `CREATE TABLE gauge(
			"name" varchar PRIMARY KEY,
			"value" double precision
			);`)

		if errCreate != nil {
			logger.Log.Error(errCreate.Error())
			tx.Rollback()
			return errCreate
		}

		logger.Log.Info("Database table named 'gauge' was been successfully created")
	}

	if !CheckTableExist(db, "counter") {
		_, errCreate := tx.ExecContext(context.TODO(), `CREATE TABLE counter(
			"name" varchar PRIMARY KEY,
			"value" bigint
			);`)

		if errCreate != nil {
			logger.Log.Error(errCreate.Error())
			tx.Rollback()
			return errCreate
		}

		logger.Log.Info("Database table named 'counter' was been successfully created")
	}

	return tx.Commit()
}

func (d *dompdb) UpdateMetricDB(mType constants.MetricType, mName string, mValue float64) error {

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
		logger.Log.Error(err.Error())
		tx.Rollback()
		return err
	}

	logger.Log.Info("Database update metric successfull.", zap.String("MetricName", mName), zap.Float64("Value", mValue))

	return tx.Commit()
}

func (d *dompdb) UpdateBatchDB(sStg storage.StorageInterface) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	gauge, counter := sStg.ReadMemStorageFields()

	tx, errBegin := d.db.Begin()

	if errBegin != nil {
		return errBegin
	}

	var err error

	for k, v := range gauge {
		q := `INSERT INTO gauge (name, value)	VALUES ($1, $2)	ON CONFLICT (name)` +
			`DO UPDATE SET name=EXCLUDED.name, value=$3;`

		_, err = tx.ExecContext(context.TODO(), q, k, v, v)

		if err != nil {
			logger.Log.Error(err.Error())
			tx.Rollback()
			return err
		}
	}

	for k, v := range counter {
		q := `INSERT INTO counter (name, value)	VALUES ($1, $2)	ON CONFLICT (name)` +
			`DO UPDATE SET name=EXCLUDED.name, value=$3;`

		_, err = tx.ExecContext(context.TODO(), q, k, v, v)

		if err != nil {
			logger.Log.Error(err.Error())
			tx.Rollback()
			return err
		}
	}

	logger.Log.Info("Updating database by metrics  batches successfully.")

	return tx.Commit()
}
