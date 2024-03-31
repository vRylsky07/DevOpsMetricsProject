package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

var Log *zap.Logger = zap.NewNop()

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

	Log = zl

	Log = zap.Must(Log, err)

	Log.Info("The Logger was successfully initialized")
	return nil
}
