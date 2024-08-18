package filesbackup

import (
	backup "DevOpsMetricsProject/internal/backups"
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	funcslib "DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
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

type FilesBackup struct {
	currentMetrics *os.File
	storeInterval  int
	log            logger.Recorder
}

func (fb *FilesBackup) IsValid() bool {
	return fb.storeInterval >= 0 && fb.currentMetrics != nil && fb.log != nil
}

func NewMetricsBackup(cfg *configs.ServerConfig, log logger.Recorder) (backup.MetricsBackup, error) {
	CreateMetricsSave(cfg.RestoreBool, log)
	tFile := CreateTempFile(cfg.TempFile, log)
	fBack := &FilesBackup{currentMetrics: tFile, storeInterval: cfg.StoreInterval, log: log}

	if !fBack.IsValid() {
		err := errors.New("creating new metrics backup failed")
		log.Error(err.Error())
		return nil, err
	}

	switch {
	case cfg.StoreInterval > 0:
		ticker := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

		go func() {
			defer ticker.Stop()
			for {
				<-ticker.C
				fBack.TransferMetricsToFile()
			}
		}()
	case cfg.StoreInterval == 0:
		log.Info("Store interval is 0, backup files will sync immediately")
	case cfg.StoreInterval < 0:
		return nil, errors.New("store interval lower then 0")
	}

	return fBack, nil
}

func (fb *FilesBackup) UpdateMetricDB(mType constants.MetricType, mName string, mValue float64) error {
	if !fb.IsValid() {
		return errors.New("UpdateMetricDB() - server is not valid")
	}

	var err error
	switch fb.storeInterval {
	case 0:
		file, errOpen := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR, 0666)
		if errOpen != nil {
			return errOpen
		}
		defer file.Close()
		err = ReplaceOrAddRowToFile(mType, mName, mValue, file, fb.log)
	default:
		err = ReplaceOrAddRowToFile(mType, mName, mValue, fb.currentMetrics, fb.log)
	}

	return err
}

func (fb *FilesBackup) TransferMetricsToFile() error {
	if !fb.IsValid() {
		return errors.New("func TransferMetricsToFile() failed, pointers are not valid")
	}

	file, err := os.Open(fb.currentMetrics.Name())

	if err != nil {
		fb.log.Error(err.Error())
		return err
	}

	defer file.Close()

	buf, errRead := io.ReadAll(file)

	if errRead != nil {
		fb.log.Error(errRead.Error())
		return errRead
	}

	savefile, errSf := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR, 0666)

	if errSf != nil {
		fb.log.Error(errSf.Error())
		return errSf
	}

	defer savefile.Close()

	errTrun := savefile.Truncate(0)
	if errTrun != nil {
		fb.log.Error(errTrun.Error())
		return errTrun
	}

	_, errSeek := savefile.Seek(0, 0)

	if errSeek != nil {
		fb.log.Error(errSeek.Error())
		return errSeek
	}

	_, errWrite := savefile.Write(buf)
	if errWrite != nil {
		fb.log.Error(errWrite.Error())
		return errWrite
	}

	fb.log.Info("Metrics was successfully transfered to save file")

	return nil
}

func (fb *FilesBackup) GetAllData() (*map[string]float64, *map[string]int) {
	errMkDir := os.MkdirAll(GetMetricsSaveFileDir(), os.ModePerm)
	if errMkDir != nil {
		fb.log.Error(errMkDir.Error())
		return nil, nil
	}

	file, err := os.OpenFile(GetMetricsSaveFilePath(), os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		fb.log.Error(err.Error())
		return nil, nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	g := make(map[string]float64)
	c := make(map[string]int)

	for scanner.Scan() {
		line := io.NopCloser(strings.NewReader(string(scanner.Bytes())))

		metricStruct, err := funcslib.DecodeMetricJSON(line)

		if err != nil {
			fb.log.Error(err.Error())
			continue
		}

		switch funcslib.ConvertStringToMetricType(metricStruct.MType) {
		case constants.GaugeType:
			if (*metricStruct).Value != nil {
				g[metricStruct.ID] = *metricStruct.Value
			}

		case constants.CounterType:
			if (*metricStruct).Delta != nil {
				c[metricStruct.ID] = int(*metricStruct.Delta)
			}
		}
	}

	return &g, &c
}

func GetMetricsSaveFilePath() string {
	return filepath.Join(".", "saved", "SavedMetrics-db.json")
}

func GetMetricsSaveFileDir() string {
	return filepath.Join(".", "saved")
}

func CreateTempFile(filename string, log logger.Recorder) *os.File {

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

	log.Info(fmt.Sprintf("\nTemporal file with current metrics was created. Path: %s", tFile.Name()))

	return tFile
}

func CreateMetricsSave(restore bool, log logger.Recorder) {
	if !restore {
		return
	}

	_, errStat := os.Stat(GetMetricsSaveFilePath())

	if errStat == nil {
		return
	}

	dir := filepath.Join(".", "saved")

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return
	}

	_, errCreate := os.Create(GetMetricsSaveFilePath())
	if errCreate != nil {
		log.Error(errCreate.Error())
		return
	}

	log.Info("New savefile was created")
}

func ReplaceOrAddRowToFile(mType constants.MetricType, mName string, mValue float64, file *os.File, log logger.Recorder) error {
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

	mJSON, errJSON := funcslib.EncodeMetricJSON(mType, mName, mValue)

	if errJSON != nil {
		log.Error(errJSON.Error())
		return errJSON
	}

	for scanner.Scan() {

		readerSave := io.NopCloser(strings.NewReader(string(scanner.Bytes())))
		mFromSave, err := funcslib.DecodeMetricJSON(readerSave)

		if err != nil {
			log.Error(err.Error())
			return err
		}

		if mFromSave.ID == mName {
			containter.Write(mJSON.Bytes())
			matched = true
		} else {
			containter.Write(scanner.Bytes())
			containter.Write([]byte("\n"))
		}
	}

	if !matched {
		containter.Write(mJSON.Bytes())
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
