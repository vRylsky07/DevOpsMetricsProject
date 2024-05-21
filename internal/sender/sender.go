package sender

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"
)

//go:generate mockgen -source=sender.go -destination=mocks/sender_mocks.go
type SenderInterface interface {
	SendMetricsHTTP(reportInterval int)
	GetStorage() storage.StorageInterface
	CreateMetricURL(mType constants.MetricType, mainURL string, name string, value float64) string
}

type dompsender struct {
	senderMemStorage storage.StorageInterface
	stopThread       bool
	cfg              *configs.ClientConfig
}

func (sStg *dompsender) IsValid() bool {
	if sStg != nil && sStg.senderMemStorage != nil && sStg.cfg != nil {
		return true
	}
	logger.Log.Error("Sender Storage is not valid")
	return false
}

func (sStg *dompsender) GetStorage() storage.StorageInterface {
	if !sStg.IsValid() {
		return nil
	}
	return sStg.senderMemStorage
}

func (sStg *dompsender) InitSenderStorage(cfg *configs.ClientConfig, newStg storage.StorageInterface) {
	sStg.senderMemStorage = newStg
	sStg.cfg = cfg
}

func (sStg *dompsender) UpdateMetrics() {
	if !sStg.IsValid() {
		return
	}

	if sStg.GetStorage() == nil {
		sStg.GetStorage().InitMemStorage()
	}

	var ticker *time.Ticker

	if sStg.cfg.PollInterval >= 0 {
		ticker = time.NewTicker(time.Duration(sStg.cfg.PollInterval) * time.Second)
		defer ticker.Stop()
	}

	for !sStg.stopThread {

		if ticker != nil {
			<-ticker.C
		}

		sStg.updateCounterMetrics()
		sStg.updateGaugeMetrics()

		if sStg.cfg.PollInterval == -1 {
			return
		}
	}
}

func (sStg *dompsender) SendMetricsHTTP() []error {
	if !sStg.IsValid() {
		return []error{errors.New("sender Storage is not valid")}
	}

	interval := sStg.cfg.ReportInterval

	var catchErrs []error

	var ticker *time.Ticker

	if interval >= 0 {
		ticker = time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
	}

	for !sStg.stopThread {

		if ticker != nil {
			<-ticker.C
		}

		if sStg == nil {
			catchErrs = append(catchErrs, errors.New("SendMetricsHTTP() FAILED! Storage of sender module is equal nil"))
			return catchErrs
		}

		gauge, counter := sStg.GetStorage().ReadMemStorageFields()

		for nameGauge, valueGauge := range gauge {
			sStg.postRequestByMetricType(true, constants.GaugeType, nameGauge, valueGauge, &catchErrs)
		}

		for nameCounter, valueCounter := range counter {
			sStg.postRequestByMetricType(true, constants.CounterType, nameCounter, float64(valueCounter), &catchErrs)
		}
		if interval == -1 {
			return catchErrs
		}
	}
	return catchErrs
}

func (sStg *dompsender) StopAgentProcessing() {
	if !sStg.IsValid() {
		return
	}
	sStg.stopThread = true
}

func CreateSender(cfg *configs.ClientConfig) *dompsender {
	senderStorage := storage.MemStorage{}
	senderStorage.InitMemStorage()

	mSender := &dompsender{}
	mSender.InitSenderStorage(cfg, &senderStorage)

	return mSender
}

func (sStg *dompsender) updateCounterMetrics() {
	if !sStg.IsValid() {
		return
	}

	if sStg.GetStorage() == nil {
		sStg.GetStorage().InitMemStorage()
	}

	sStg.GetStorage().UpdateMetricByName(constants.AddOperation, constants.CounterType, "PollCount", 1)
}

func (sStg *dompsender) updateGaugeMetrics() {
	if !sStg.IsValid() {
		return
	}

	if sStg.GetStorage() == nil {
		sStg.GetStorage().InitMemStorage()
	}

	mFromRuntime := &runtime.MemStats{}
	runtime.ReadMemStats(mFromRuntime)

	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "RandomValue", functionslibrary.GenerateRandomValue(-10000, 10000, 3))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "Alloc", float64(mFromRuntime.Alloc))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "BuckHashSys", float64(mFromRuntime.BuckHashSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "Frees", float64(mFromRuntime.Frees))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "GCCPUFraction", float64(mFromRuntime.GCCPUFraction))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "GCSys", float64(mFromRuntime.GCSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "HeapAlloc", float64(mFromRuntime.HeapAlloc))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "HeapIdle", float64(mFromRuntime.HeapIdle))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "HeapInuse", float64(mFromRuntime.HeapInuse))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "HeapReleased", float64(mFromRuntime.HeapReleased))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "HeapObjects", float64(mFromRuntime.HeapObjects))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "HeapSys", float64(mFromRuntime.HeapSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "LastGC", float64(mFromRuntime.LastGC))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "Lookups", float64(mFromRuntime.Lookups))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "MCacheInuse", float64(mFromRuntime.MCacheInuse))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "MCacheSys", float64(mFromRuntime.MCacheSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "MSpanInuse", float64(mFromRuntime.MSpanInuse))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "MSpanSys", float64(mFromRuntime.MSpanSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "Mallocs", float64(mFromRuntime.Mallocs))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "NextGC", float64(mFromRuntime.NextGC))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "NumForcedGC", float64(mFromRuntime.NumForcedGC))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "NumGC", float64(mFromRuntime.NumGC))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "OtherSys", float64(mFromRuntime.OtherSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "PauseTotalNs", float64(mFromRuntime.PauseTotalNs))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "StackInuse", float64(mFromRuntime.StackInuse))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "StackSys", float64(mFromRuntime.StackSys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "Sys", float64(mFromRuntime.Sys))
	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "TotalAlloc", float64(mFromRuntime.TotalAlloc))
}

func (sStg *dompsender) postRequestByMetricType(compress bool, mType constants.MetricType, mName string, mValue float64, catchErrs *[]error) {

	if !sStg.IsValid() {
		return
	}

	sendURL := "http://" + sStg.cfg.Address + "/update/"

	mJSON, jsonErr := functionslibrary.EncodeMetricJSON(mType, mName, mValue)

	if jsonErr != nil {
		*catchErrs = append(*catchErrs, jsonErr)
		logger.Log.Error(jsonErr.Error())
		return
	}

	if compress {
		zipped, compErr := functionslibrary.CompressData(mJSON.Bytes())
		if compErr == nil {
			mJSON = zipped
		}
	}

	client := http.Client{}

	req, errReq := http.NewRequest("POST", sendURL, mJSON)

	if errReq != nil {
		logger.Log.Error(errReq.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")

	if compress {
		req.Header.Add("Content-Encoding", "gzip ")
	}

	resp, errDo := client.Do(req)

	if errDo != nil {
		errStr := "Server is not responding. URL to send was: " + sendURL
		*catchErrs = append(*catchErrs, errors.New(errStr))
		logger.Log.Error(errStr)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Info(fmt.Sprintf(`Metric %s update failed! Status code: %d`, mName, resp.StatusCode))
		return
	}

	logger.Log.Info(fmt.Sprintf(`Metric %s update was successful! Status code: %d`, mName, resp.StatusCode))
}
