package server

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/logger"
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

	row := db.QueryRowContext(context.Background(),
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

	if !CheckTableExist(db, "gauge") {
		_, errCreate := db.ExecContext(context.TODO(), `CREATE TABLE gauge(
			"name" varchar PRIMARY KEY,
			"value" double precision
			);`)

		if errCreate != nil {
			logger.Log.Error(errCreate.Error())
			return errCreate
		}

		logger.Log.Info("Database table named 'gauge' was been successfully created")
	}

	if !CheckTableExist(db, "counter") {
		_, errCreate := db.ExecContext(context.TODO(), `CREATE TABLE counter(
			"name" varchar PRIMARY KEY,
			"value" int
			);`)

		if errCreate != nil {
			logger.Log.Error(errCreate.Error())
			return errCreate
		}

		logger.Log.Info("Database table named 'counter' was been successfully created")
	}

	return nil

}

func UpdateMetricDB(db *sql.DB, mType constants.MetricType, mName string, mValue float64) error {

	q := fmt.Sprintf(`INSERT INTO %s (name, value)
	VALUES ($1, $2)
	ON CONFLICT (name)
	DO UPDATE SET
	name=EXCLUDED.name,
	value=$3;`, functionslibrary.ConvertMetricTypeToString(mType))

	_, err := db.ExecContext(context.TODO(), q, mName, mValue, mValue)

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	logger.Log.Info("Database update metric successfull.", zap.String("MetricName", mName), zap.Float64("Value", mValue))

	return nil
}
