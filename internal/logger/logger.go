package logger

import (
	"bufio"
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

func (dl *dompLog) GetLine() {

	file, err := os.OpenFile("logs/server_domp.log", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}

	var scanner *bufio.Scanner = bufio.NewScanner(file)
	if !scanner.Scan() {
		return
	}

	Log.Info(scanner.Text())
}
