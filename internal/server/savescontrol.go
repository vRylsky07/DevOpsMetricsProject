package server

import (
	"DevOpsMetricsProject/internal/configs"
	funcslib "DevOpsMetricsProject/internal/funcslib"
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
		serv.log.Info("StoreInterval is 0 and metrics will save after update immediately")
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
		serv.log.Error(err.Error())
		return err
	}

	defer file.Close()

	buf, errRead := io.ReadAll(file)

	if errRead != nil {
		serv.log.Error(errRead.Error())
		return errRead
	}

	savefile, errSf := os.OpenFile(serv.savefile.Savefile.Name(), os.O_RDWR, 0666)

	if errSf != nil {
		serv.log.Error(errSf.Error())
		return errSf
	}

	errTrun := savefile.Truncate(0)
	if errTrun != nil {
		serv.log.Error(errTrun.Error())
		return errTrun
	}

	_, errSeek := savefile.Seek(0, 0)

	if errSeek != nil {
		serv.log.Error(errSeek.Error())
		return errSeek
	}

	_, errWrite := savefile.Write(buf)
	if errWrite != nil {
		serv.log.Error(errWrite.Error())
		return errWrite
	}

	serv.log.Info("Metrics was successfully transfered to save file")

	return nil
}

func (serv *dompserver) SaveCurrentMetrics(b *bytes.Buffer) error {
	if !serv.IsValid() {
		return errors.New("SaveCurrentMetrics() - server is not valid")
	}

	var err error
	switch serv.savefile.StoreInterval {
	case 0:
		err = ReplaceOrAddRowToFile(serv.savefile.Savefile, b, serv.log)
	default:
		err = ReplaceOrAddRowToFile(serv.currentMetrics, b, serv.log)
	}

	return err
}

func GetMetricsSaveFilePath() string {
	return filepath.Join(".", "saved", "SavedMetrics-db.json")
}

func GetMetricsSaveFileDir() string {
	return filepath.Join(".", "saved")
}

func CreateTempFile(filename string, restore bool, log logger.Recorder) *os.File {

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
			log.Error(err.Error())
			return nil
		}

		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Error(err.Error())
			return nil
		}
	}

	tFile, errCreate := os.CreateTemp(dir, "*_"+name)
	if errCreate != nil {
		log.Error(errCreate.Error())
		return nil
	}

	if restore {
		errMkDir := os.MkdirAll(GetMetricsSaveFileDir(), os.ModePerm)

		if errMkDir != nil {
			log.Error(errMkDir.Error())
			return nil
		}

		file, err := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {
			log.Error(err.Error())
			return nil
		}

		defer file.Close()

		buf, errRead := io.ReadAll(file)
		if errRead == nil {
			tFile.Write(buf)
		}
	}

	log.Info(fmt.Sprintf("\nTemporal file with current metrics was created. Path: %s", tFile.Name()))

	return tFile
}

func CreateMetricsSave(interval int, log logger.Recorder) *MetricsSave {
	dir := filepath.Join(".", "saved")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	file, errCreate := os.Create(GetMetricsSaveFilePath())
	if errCreate != nil {
		log.Error(errCreate.Error())
		return nil
	}

	log.Info("New savefile was created")
	return &MetricsSave{interval, file}
}

func RestoreData(cfg *configs.ServerConfig, sStg storage.MetricsRepository, log logger.Recorder) *MetricsSave {

	if !cfg.RestoreBool || sStg == nil {
		return CreateMetricsSave(cfg.StoreInterval, log)
	}

	errMkDir := os.MkdirAll(GetMetricsSaveFileDir(), os.ModePerm)
	if errMkDir != nil {
		log.Error(errMkDir.Error())
		return nil
	}

	file, err := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Error(err.Error())
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := io.NopCloser(strings.NewReader(string(scanner.Bytes())))

		metricStruct, err := funcslib.DecodeMetricJSON(line)

		if err != nil {
			log.Error(err.Error())
			continue
		}

		errUpdate := funcslib.UpdateStorageInterfaceByMetricStruct(sStg, metricStruct)
		if errUpdate != nil {
			log.Error(errUpdate.Error())
			continue
		}
	}

	log.Info("Storage was successfully restored from save file")
	return &MetricsSave{cfg.StoreInterval, file}
}

func ReplaceOrAddRowToFile(file *os.File, b *bytes.Buffer, log logger.Recorder) error {
	openF, err := os.OpenFile(file.Name(), os.O_RDWR, 0666)

	if err != nil {
		log.Error(err.Error())
		return err
	}

	defer openF.Close()

	scanner := bufio.NewScanner(openF)

	var newBuf []byte
	containter := bytes.NewBuffer(newBuf)

	matched := false

	readerBuf := io.NopCloser(strings.NewReader(b.String()))
	mFromBuf, errBuf := funcslib.DecodeMetricJSON(readerBuf)

	if errBuf != nil {
		log.Error(errBuf.Error())
		return errBuf
	}

	for scanner.Scan() {

		readerSave := io.NopCloser(strings.NewReader(string(scanner.Bytes())))
		mFromSave, err := funcslib.DecodeMetricJSON(readerSave)

		if err != nil {
			log.Error(err.Error())
			return err
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
		log.Error(errTrun.Error())
		return errTrun
	}

	_, errSeek := file.Seek(0, 0)

	if errSeek != nil {
		log.Error(errSeek.Error())
		return errSeek
	}

	file.Write(containter.Bytes())
	return nil
}
