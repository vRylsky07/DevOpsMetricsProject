package sender

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/funcslib"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"go.uber.org/zap"
)

//go:generate mockgen -source=sender.go -destination=mocks/sender_mocks.go
type MetricsProvider interface {
	SendMetricsHTTP(reportInterval int)
	GetStorage() storage.MetricsRepository
	CreateMetricURL(mType constants.MetricType, mainURL string, name string, value float64) string
}

type dompsender struct {
	senderMemStorage storage.MetricsRepository
	stopThread       bool
	cfg              *configs.ClientConfig
	log              logger.Recorder
}

func (sStg *dompsender) GetLogger() logger.Recorder {
	return sStg.log
}

func (sStg *dompsender) IsValid() bool {
	if sStg != nil && sStg.senderMemStorage != nil && sStg.cfg != nil {
		return true
	}
	sStg.log.Error("Sender Storage is not valid")
	return false
}

func (sStg *dompsender) GetStorage() storage.MetricsRepository {
	if !sStg.IsValid() {
		return nil
	}
	return sStg.senderMemStorage
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

	var catchErrs []error

	var ticker *time.Ticker

	if sStg.cfg.ReportInterval >= 0 {
		ticker = time.NewTicker(time.Duration(sStg.cfg.ReportInterval) * time.Second)
		defer ticker.Stop()
	}

	for !sStg.stopThread {

		if ticker != nil {
			<-ticker.C
		}

		sStg.ManageRequests(&catchErrs, ticker)

		if sStg.cfg.ReportInterval == -1 {
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

func CreateSender(cfg *configs.ClientConfig) (*dompsender, error) {
	senderStorage := storage.MemStorage{}
	senderStorage.InitMemStorage()

	log, err := logger.Initialize(cfg.Loglevel, "agent_")

	if err != nil {
		return nil, err
	}

	mSender := &dompsender{senderMemStorage: &senderStorage, cfg: cfg, log: log}
	return mSender, nil
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

	sStg.GetStorage().UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "RandomValue", funcslib.GenerateRandomValue(-10000, 10000, 3))
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

func (sStg *dompsender) postRequestByMetricType(ticker *time.Ticker, mName string, mJSON *bytes.Buffer, encErr error, catchErrs *[]error) {
	if ticker != nil {
		ticker.Stop()
		defer ticker.Reset(time.Duration(sStg.cfg.ReportInterval) * time.Second)
	}

	if !sStg.IsValid() {
		return
	}

	if encErr != nil {
		sStg.log.Error(encErr.Error())
		*catchErrs = append(*catchErrs, encErr)
		return
	}

	batchStr := ""

	if sStg.cfg.UseBatches {
		batchStr = "s"
	}

	sendURL := "http://" + sStg.cfg.Address + "/update" + batchStr + "/"

	if sStg.cfg.CompressData {
		zipped, compErr := funcslib.CompressData(mJSON.Bytes())
		if compErr == nil {
			mJSON = zipped
		} else {
			sStg.log.Error(compErr.Error())
		}
	}

	client := http.Client{}

	req, errReq := http.NewRequest("POST", sendURL, mJSON)

	if errReq != nil {
		sStg.log.Error(errReq.Error())
		return
	}

	req.Header.Add("Content-Type", "application/json")

	if sStg.cfg.CompressData {
		req.Header.Add("Content-Encoding", "gzip ")
	}

	var resp *http.Response
	var errDo error

	for _, v := range *constants.GetRetryIntervals() {
		if v != 0 {
			sStg.log.Info("Server is not responding. Retry to do post request...")
			timer := time.NewTimer(time.Duration(v) * time.Second)
			<-timer.C
		}
		resp, errDo = client.Do(req)

		if errDo == nil {
			defer resp.Body.Close()
			break
		}
	}

	if errDo != nil {
		errStr := "Server is not responding. URL to send was: " + sendURL
		*catchErrs = append(*catchErrs, errors.New(errStr))
		sStg.log.Error(errStr)
		return
	}

	if resp.StatusCode != http.StatusOK {
		sStg.log.Info(fmt.Sprintf(`Metrics update failed! Status code: %d`, resp.StatusCode))
		return
	}

	sStg.log.Info(fmt.Sprintf(`Metric%s update was successful! Status code: %d`, batchStr, resp.StatusCode), zap.String("MetricName", mName))
}

func (sStg *dompsender) ManageRequests(catchErrs *[]error, ticker *time.Ticker) {
	var mJSON *bytes.Buffer
	var errJSON error

	switch sStg.cfg.UseBatches {
	case true:
		mJSON, errJSON = funcslib.EncodeBatchJSON(sStg.GetStorage())
		sStg.postRequestByMetricType(ticker, "batch", mJSON, errJSON, catchErrs)
	case false:

		gauge, counter := sStg.GetStorage().ReadMemStorageFields()

		for nameGauge, valueGauge := range gauge {
			mJSON, errJSON = funcslib.EncodeMetricJSON(constants.GaugeType, nameGauge, valueGauge)
			sStg.postRequestByMetricType(ticker, nameGauge, mJSON, errJSON, catchErrs)
		}

		for nameCounter, valueCounter := range counter {
			mJSON, errJSON = funcslib.EncodeMetricJSON(constants.CounterType, nameCounter, float64(valueCounter))
			sStg.postRequestByMetricType(ticker, nameCounter, mJSON, errJSON, catchErrs)
		}
	}
}
