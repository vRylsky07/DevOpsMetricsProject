package server

import (
	"DevOpsMetricsProject/internal/constants"
	funcslib "DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"context"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func RunDB(dsn string) (*sql.DB, error) {
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

	logger.Log.Info("Connection to Database is OK")
	return db, nil
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

func UpdateMetricDB(db *sql.DB, mType constants.MetricType, mName string, mValue float64) error {

	tx, errBegin := db.Begin()

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

func UpdateBatchDB(db *sql.DB, sStg storage.StorageInterface) error {
	gauge, counter := sStg.ReadMemStorageFields()

	tx, errBegin := db.Begin()

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
