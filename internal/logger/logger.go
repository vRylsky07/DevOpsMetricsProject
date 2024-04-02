package logger

import (
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type dompLog struct {
	*zap.Logger
}

func (dl *dompLog) ErrorHTTP(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	dl.Error(err.Error())
}

var Log *dompLog = &dompLog{zap.NewNop()}

func Initialize(level string, filePrefix string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl

	outputPath := filepath.Join(".", "logs")
	os.MkdirAll(outputPath, os.ModePerm)

	outputPath = filepath.Join(outputPath, filePrefix+"domp.log")

	cfg.OutputPaths = []string{outputPath, "stderr"}

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = &dompLog{zl}

	zapLogger := zap.Must(Log.Logger, err)

	Log = &dompLog{zapLogger}

	Log.Info("The Logger was successfully initialized")
	return nil
}
