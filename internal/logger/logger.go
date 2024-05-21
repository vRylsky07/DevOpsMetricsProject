package logger

import (
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerInterface interface {
	Error(msg string, fields ...interface{})
	ErrorHTTP(w http.ResponseWriter, err error, code int)
	Info(msg string, fields ...interface{})
}

type dompLogZap struct {
	*zap.Logger
}

func (dl *dompLogZap) Error(msg string, fields ...interface{}) {
	if len(fields) <= 0 {
		dl.Logger.Error(msg)
		return
	}

	container := []zapcore.Field{}

	for _, field := range fields {
		f, ok := field.(zapcore.Field)
		if ok {
			container = append(container, f)
		} else {
			dl.Logger.Error("zapcore.Field convertation failed")
		}
	}

	dl.Logger.Error(msg, container...)
}

func (dl *dompLogZap) Info(msg string, fields ...interface{}) {
	if len(fields) <= 0 {
		dl.Logger.Info(msg)
		return
	}

	container := []zapcore.Field{}

	for _, field := range fields {
		f, ok := field.(zapcore.Field)
		if ok {
			container = append(container, f)
		} else {
			dl.Logger.Error("zapcore.Field convertation failed")
		}
	}

	dl.Logger.Info(msg, container...)
}

func (dl *dompLogZap) ErrorHTTP(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
	dl.Error(err.Error())
}

var Log LoggerInterface = &dompLogZap{zap.NewNop()}

func Initialize(level string, filePrefix string) error {

	var initLog *dompLogZap

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

	initLog = &dompLogZap{zl}

	zapLogger := zap.Must(initLog.Logger, err)

	Log = &dompLogZap{zapLogger}

	Log.Info("The Logger was successfully initialized")
	return nil
}
