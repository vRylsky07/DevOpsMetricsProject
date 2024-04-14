package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type dompserver struct {
	coreMux        *chi.Mux
	coreStg        storage.StorageInterface
	currentMetrics *os.File
	cfg            *configs.ServerConfig
	savefile       *MetricsSave
}

func (serv *dompserver) IsValid() bool {
	return serv.coreMux != nil || serv.coreStg != nil || serv.currentMetrics != nil || serv.cfg != nil || serv.savefile != nil || serv.savefile.IsValid()
}

func (serv *dompserver) TransferMetricsToFile() error {
	file, err := os.Open(serv.currentMetrics.Name())

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	buf, errRead := io.ReadAll(file)

	if errRead != nil {
		logger.Log.Error(errRead.Error())
		return errRead
	}
	serv.savefile.Savefile.Truncate(0)
	serv.savefile.Savefile.Seek(0, 0)
	_, errWrite := serv.savefile.Savefile.Write(buf)
	if errWrite != nil {
		logger.Log.Error(errWrite.Error())
		return errWrite
	}

	logger.Log.Info("Metrics was succesfully transfered to save file")

	return nil
}

func (serv *dompserver) SaveCurrentMetrics(b *bytes.Buffer) {
	switch serv.savefile.StoreInterval {
	case 0:
		serv.savefile.Savefile.Write(b.Bytes())
	default:
		serv.currentMetrics.Write(b.Bytes())
	}
}

func (serv *dompserver) StartSaveMetricsThread() {
	if serv.savefile.StoreInterval == 0 {
		logger.Log.Info("StoreInterval is 0 and metrics will save after update immediately")
		return
	}

	if serv.savefile != nil && serv.savefile.StoreInterval > 0 && serv.savefile.Savefile != nil {
		go func() {
			for {
				time.Sleep(time.Duration(serv.savefile.StoreInterval) * time.Second)
				serv.TransferMetricsToFile()
			}
		}()
	}

}

func Start() {
	dompserv := CreateNewServer()
	dompserv.StartSaveMetricsThread()
	if !dompserv.IsValid() {
		logger.Log.Info(
			"Server initialization failed!  ",
			zap.Bool("coreMux is nil?", (dompserv.coreMux == nil)),
			zap.Bool("coreStg is nil?", (dompserv.coreStg == nil)),
			zap.Bool("currentMetrics is nil?", (dompserv.currentMetrics == nil)),
			zap.Bool("cfg is nil?", (dompserv.cfg == nil)),
			zap.Bool("savefile is not valid?", (!dompserv.savefile.IsValid())),
		)
		return
	}
	logger.Log.Info("Server was successfully initialized!")
	err := http.ListenAndServe(dompserv.cfg.Address, dompserv.coreMux)
	if err != nil {
		panic(err)
	}

}

func CreateNewServer() *dompserver {
	dompserv := NewDompServer()

	dompserv.coreMux.Use(dompserv.WithResponseLog)
	dompserv.coreMux.Use(dompserv.WithRequestLog)
	dompserv.coreMux.Use(gzipHandle)
	dompserv.coreMux.Use(DecompressHandler)

	dompserv.coreMux.Get("/", dompserv.GetMainPageHandler)
	dompserv.coreMux.Route("/update", func(r chi.Router) {
		r.Post("/", dompserv.MetricHandlerJSON)
		r.Get("/", dompserv.IncorrectRequestHandler)
		r.Post("/{mType}/{mName}/{mValue}", dompserv.UpdateMetricHandler)
		r.Get("/{mType}/{mName}/{mValue}", dompserv.IncorrectRequestHandler)
	})
	dompserv.coreMux.Route("/value", func(r chi.Router) {
		r.Post("/", dompserv.MetricHandlerJSON)
		r.Get("/{mType}/{mName}", dompserv.GetMetricHandler)
	})
	return dompserv
}

func NewDompServer() *dompserver {
	coreMux := chi.NewRouter()
	coreStg := &storage.MemStorage{}
	coreStg.InitMemStorage()
	cfg := configs.CreateServerConfig()
	logger.Initialize(cfg.Loglevel, "server_")

	savefile := RestoreData(cfg, coreStg, GetMetricsSaveFilePath())

	if !cfg.RestoreBool || savefile == nil {
		savefile = CreateMetricsSave(cfg.StoreInterval)
	}

	serv := &dompserver{
		coreMux:        coreMux,
		coreStg:        coreStg,
		currentMetrics: CreateTempFile(cfg.TempFile),
		cfg:            cfg,
		savefile:       savefile,
	}
	return serv
}

func CreateTempFile(filename string) *os.File {
	dir := filepath.Join(os.TempDir(), "domp_temp")

	err := os.RemoveAll(dir)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}

	tFile, errCreate := os.CreateTemp(dir, "*_"+filename)
	if errCreate != nil {
		logger.Log.Error(errCreate.Error())
		return nil
	}

	logger.Log.Info(fmt.Sprintf("\nTemporal file with current metrics was created. Path: %s", dir))

	return tFile
}

func GetMetricsSaveFilePath() string {
	return filepath.Join(".", "saved", "SavedMetrics-db.json")
}

func CreateMetricsSave(interval int) *MetricsSave {
	if interval < 0 {
		logger.Log.Error("CreateSavingThread() failed. Store interval value cannot be negative")
		return nil
	}

	dir := filepath.Join(".", "saved")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}

	file, errCreate := os.Create(GetMetricsSaveFilePath())
	if errCreate != nil {
		logger.Log.Error(errCreate.Error())
		return nil
	}

	return &MetricsSave{interval, file}
}

func RestoreData(cfg *configs.ServerConfig, sStg storage.StorageInterface, path string) *MetricsSave {
	if !cfg.RestoreBool || sStg == nil {
		logger.Log.Info("Restore data skipped")
		return nil
	}

	file, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := io.NopCloser(strings.NewReader(string(scanner.Bytes())))

		metricStruct, err := functionslibrary.DecodeMetricJSON(line)

		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}

		errUpdate := functionslibrary.UpdateStorageInterfaceByMetricStruct(sStg, functionslibrary.ConvertStringToMetricType(metricStruct.MType), metricStruct)
		if errUpdate != nil {
			logger.Log.Error(errUpdate.Error())
			continue
		}
	}

	logger.Log.Info("Storage was successfully restored from save file")
	return &MetricsSave{cfg.StoreInterval, file}
}

//конструкторы NewRouter
//логгер на два
//конфиг передвинуть
