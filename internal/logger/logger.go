package logger

import (
	"fmt"
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

	//var fieldsArr []zapcore.Field = make([]zapcore.Field, len(fields))

	for _, field := range fields {
		f, ok := field.(zapcore.Field)
		if ok {
			fmt.Printf("IS OK %s", f.Key)
			//fieldsArr = append(fieldsArr, f)
		} else {
			fmt.Println("NOT OK")
			return
		}
	}
	fmt.Println("WOOOPS")
	//dl.Logger.Error(msg, fieldsArr...)
}

func (dl *dompLogZap) Info(msg string, fields ...interface{}) {
	if len(fields) <= 0 {
		dl.Logger.Info(msg)
		return
	}

	fmt.Println("START FIELD CHECK")
	//var fieldsArr []zapcore.Field = make([]zapcore.Field, len(fields))

	for _, field := range fields {
		f, ok := field.(zapcore.Field)
		if ok {
			fmt.Printf("IS OK %v", f)

			//zF := zapcore.Field{Key: f.Key, Type: f.Type, Integer: f.Integer, String: f.String, Interface: f.Interface}
			//fieldsArr = append(fieldsArr, zF)
		} else {
			fmt.Println("NOT OK")
			return
		}
	}

	fmt.Println("WOOOPS")
	dl.Logger.Info(msg, FiledInt("KEY", 77))
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

func FiledInt(k string, v int) zapcore.Field {
	return zap.Int64(k, int64(v))
}

func FieldString(k string, v string) zapcore.Field {
	return zap.String(k, v)
}
