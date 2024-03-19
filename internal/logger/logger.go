package logger

import "go.uber.org/zap"

var Log *zap.Logger = zap.NewNop()

func Initialize(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.OutputPaths = []string{"./logs/domp.log", "stderr"}

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl

	Log = zap.Must(Log, err)

	Log.Info("The Logger was successfully initialized")
	return nil
}
