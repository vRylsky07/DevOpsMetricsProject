package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	functionslibrary "DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (serv *dompserver) StartSaveMetricsThread() {
	if !serv.IsValid() {
		return
	}

	if serv.savefile.StoreInterval == 0 {
		logger.Log.Info("StoreInterval is 0 and metrics will save after update immediately")
		return
	}

	if serv.savefile != nil && serv.savefile.StoreInterval > 0 && serv.savefile.Savefile != nil {

		ticker := time.NewTicker(time.Duration(serv.savefile.StoreInterval) * time.Second)

		go func() {
			defer ticker.Stop()
			for {
				<-ticker.C
				serv.TransferMetricsToFile()
			}
		}()
	}
}

func (serv *dompserver) TransferMetricsToFile() error {
	if !serv.IsValid() {
		return errors.New("DOMP Server is not valid")
	}

	file, err := os.Open(serv.currentMetrics.Name())

	if err != nil {
		logger.Log.Error(err.Error())
		return err
	}

	defer file.Close()

	buf, errRead := io.ReadAll(file)

	if errRead != nil {
		logger.Log.Error(errRead.Error())
		return errRead
	}

	savefile, errSf := os.OpenFile(serv.savefile.Savefile.Name(), os.O_RDWR, 0666)

	if errSf != nil {
		logger.Log.Error(errSf.Error())
		return errSf
	}

	errTrun := savefile.Truncate(0)
	if errTrun != nil {
		logger.Log.Error(errTrun.Error())
		return errTrun
	}

	_, errSeek := savefile.Seek(0, 0)

	if errSeek != nil {
		logger.Log.Error(errSeek.Error())
		return errSeek
	}

	_, errWrite := savefile.Write(buf)
	if errWrite != nil {
		logger.Log.Error(errWrite.Error())
		return errWrite
	}

	logger.Log.Info("Metrics was successfully transfered to save file")

	return nil
}

func (serv *dompserver) SaveCurrentMetrics(b *bytes.Buffer) {
	if !serv.IsValid() {
		return
	}
	switch serv.savefile.StoreInterval {
	case 0:
		ReplaceOrAddRowToFile(serv.savefile.Savefile, b)
	default:
		ReplaceOrAddRowToFile(serv.currentMetrics, b)
	}
}

func GetMetricsSaveFilePath() string {
	return filepath.Join(".", "saved", "SavedMetrics-db.json")
}

func GetMetricsSaveFileDir() string {
	return filepath.Join(".", "saved")
}

func CreateTempFile(filename string, restore bool) *os.File {

	noSepStr := strings.Split(filename, "/")

	dir := os.TempDir()

	skipMkDir := true

	for i, str := range noSepStr {
		if i == len(noSepStr)-1 {
			break
		}
		dir = filepath.Join(dir, str)
		skipMkDir = false
	}

	name := noSepStr[len(noSepStr)-1]

	if !skipMkDir {
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
	}

	tFile, errCreate := os.CreateTemp(dir, "*_"+name)
	if errCreate != nil {
		logger.Log.Error(errCreate.Error())
		return nil
	}

	if restore {
		errMkDir := os.MkdirAll(GetMetricsSaveFileDir(), os.ModePerm)

		if errMkDir != nil {
			logger.Log.Error(errMkDir.Error())
			return nil
		}

		file, err := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			logger.Log.Error(err.Error())
			return nil
		}

		defer file.Close()

		buf, errRead := io.ReadAll(file)
		if errRead == nil {
			tFile.Write(buf)
		}
	}

	logger.Log.Info(fmt.Sprintf("\nTemporal file with current metrics was created. Path: %s", tFile.Name()))

	return tFile
}

func CreateMetricsSave(interval int) *MetricsSave {

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

	logger.Log.Info("New savefile was created")
	return &MetricsSave{interval, file}
}

func RestoreData(cfg *configs.ServerConfig, sStg storage.StorageInterface) *MetricsSave {

	if cfg.SaveMode != constants.FileMode {
		return nil
	}

	if !cfg.RestoreBool || sStg == nil {
		logger.Log.Info("Restore data skipped")
		return CreateMetricsSave(cfg.StoreInterval)
	}

	errMkDir := os.MkdirAll(GetMetricsSaveFileDir(), os.ModePerm)
	if errMkDir != nil {
		logger.Log.Error(errMkDir.Error())
		return nil
	}

	file, err := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}

	defer file.Close()

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

func ReplaceOrAddRowToFile(file *os.File, b *bytes.Buffer) {
	openF, err := os.OpenFile(file.Name(), os.O_RDWR, 0666)

	if err != nil {
		logger.Log.Error(err.Error())
		return
	}

	defer openF.Close()

	scanner := bufio.NewScanner(openF)

	var newBuf []byte
	containter := bytes.NewBuffer(newBuf)

	matched := false

	readerBuf := io.NopCloser(strings.NewReader(b.String()))
	mFromBuf, errBuf := functionslibrary.DecodeMetricJSON(readerBuf)

	if errBuf != nil {
		logger.Log.Error(errBuf.Error())
		return
	}

	for scanner.Scan() {

		readerSave := io.NopCloser(strings.NewReader(string(scanner.Bytes())))
		mFromSave, err := functionslibrary.DecodeMetricJSON(readerSave)

		if err != nil {
			logger.Log.Error(err.Error())
			return
		}

		if mFromSave.ID == mFromBuf.ID {
			containter.Write(b.Bytes())
			matched = true
		} else {
			containter.Write(scanner.Bytes())
			containter.Write([]byte("\n"))
		}
	}

	if !matched {
		containter.Write(b.Bytes())
	}

	errTrun := file.Truncate(0)
	if errTrun != nil {
		logger.Log.Error(errTrun.Error())
		return
	}

	_, errSeek := file.Seek(0, 0)

	if errSeek != nil {
		logger.Log.Error(errSeek.Error())
		return
	}

	file.Write(containter.Bytes())
}
