package server

import (
	"DevOpsMetricsProject/internal/logger"
	"context"
	"database/sql"
)

func RunDB(dsn string) (*sql.DB, error) {
	logger.Log.Info("DB: " + dsn)
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
