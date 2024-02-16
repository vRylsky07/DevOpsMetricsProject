package sender

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/functionslibrary"
	"DevOpsMetricsProject/internal/storage"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

//go:generate mockgen -source=sender.go -destination=mocks/sender_mocks.go
type SenderInterface interface {
	SendMetricsHTTP(reportInterval int)
	GetStorage() storage.StorageInterface
	CreateMetricURL(mType constants.MetricType, mainURL string, name string, value float64) string
}

type SenderStorage struct {
	senderMemStorage storage.StorageInterface
	stopThread       bool
}

func (sStg *SenderStorage) GetStorage() storage.StorageInterface {
	return sStg.senderMemStorage
}

func (sStg *SenderStorage) InitSenderStorage(newStg storage.StorageInterface) {
	sStg.senderMemStorage = newStg
	sStg.stopThread = false
}

// обновление метрик
func (sStg *SenderStorage) UpdateMetrics(pollInterval int) {

	if sStg.GetStorage() == nil {
		sStg.GetStorage().InitMemStorage()
	}

	for !sStg.stopThread {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		sStg.updateCounterMetrics()
		sStg.updateGaugeMetrics()
		if pollInterval == -1 {
			return
		}
	}
}

// отправка метрик
func (sStg *SenderStorage) SendMetricsHTTP(reportInterval int) []error {

	var catchErrs []error

	for !sStg.stopThread {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		if sStg == nil {
			catchErrs = append(catchErrs, errors.New("SendMetricsHTTP() FAILED! Storage of sender module is equal nil"))
			return catchErrs
		}

		gauge, counter := sStg.GetStorage().ReadMemStorageFields()

		for nameGauge, valueGauge := range gauge {
			finalURL := sStg.CreateMetricURL(constants.GaugeType, "http://localhost:8080", nameGauge, valueGauge)
			resp, err := http.Post(finalURL, "text/plain", nil)
			if err != nil {
				catchErrs = append(catchErrs, errors.New("Server is not responding. URL to send was: "+finalURL))
				continue
			}
			defer resp.Body.Close()
			fmt.Println(finalURL, "Gauge update was successful! Status code: ", resp.StatusCode)
		}

		for nameCounter, valueCounter := range counter {
			finalURL := sStg.CreateMetricURL(constants.CounterType, "http://localhost:8080", nameCounter, float64(valueCounter))
			resp, err := http.Post(finalURL, "text/plain", nil)
			if err != nil {
				catchErrs = append(catchErrs, errors.New("Server is not responding. URL to send was: "+finalURL))
				continue
			}
			defer resp.Body.Close()
			fmt.Println(finalURL, " Counter update was successful! Status code: ", resp.StatusCode)
		}
		if reportInterval == -1 {
			return catchErrs
		}
	}
	return catchErrs
}

func (sStg *SenderStorage) StopAgentProcessing() {
	sStg.stopThread = true
}

// конкатенация URL на основе данных метрики
func (sStg *SenderStorage) CreateMetricURL(mType constants.MetricType, mainURL string, name string, value float64) string {
	mTypeString := ""

	switch mType {
	case constants.GaugeType:
		mTypeString = "/gauge/"
	case constants.CounterType:
		mTypeString = "/counter/"
	}
	return mainURL + "/update" + mTypeString + name + "/" + strconv.FormatFloat(value, 'f', 2, 64)
}

func CreateSender() *SenderStorage {
	senderStorage := storage.MemStorage{}
	senderStorage.InitMemStorage()

	mSender := &SenderStorage{}
	mSender.InitSenderStorage(&senderStorage)

	return mSender
}

func (sStg *SenderStorage) updateCounterMetrics() {

	if sStg.GetStorage() == nil {
		sStg.GetStorage().InitMemStorage()
	}

	sStg.GetStorage().UpdateMetricByName(constants.CounterType, "PollCount", 1)
}

func (sStg *SenderStorage) updateGaugeMetrics() {

	if sStg.GetStorage() == nil {
		sStg.GetStorage().InitMemStorage()
	}

	mFromRuntime := &runtime.MemStats{}
	runtime.ReadMemStats(mFromRuntime)

	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "RandomValue", functionslibrary.GenerateRandomValue(-10000, 10000, 3))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "Alloc", float64(mFromRuntime.Alloc))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "BuckHashSys", float64(mFromRuntime.BuckHashSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "Frees", float64(mFromRuntime.Frees))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "GCCPUFraction", float64(mFromRuntime.GCCPUFraction))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "GCSys", float64(mFromRuntime.GCSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "HeapAlloc", float64(mFromRuntime.HeapAlloc))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "HeapIdle", float64(mFromRuntime.HeapIdle))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "HeapInuse", float64(mFromRuntime.HeapInuse))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "HeapReleased", float64(mFromRuntime.HeapReleased))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "HeapObjects", float64(mFromRuntime.HeapObjects))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "HeapSys", float64(mFromRuntime.HeapSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "LastGC", float64(mFromRuntime.LastGC))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "Lookups", float64(mFromRuntime.Lookups))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "MCacheInuse", float64(mFromRuntime.MCacheInuse))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "MCacheSys", float64(mFromRuntime.MCacheSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "MSpanInuse", float64(mFromRuntime.MSpanInuse))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "MSpanSys", float64(mFromRuntime.MSpanSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "Mallocs", float64(mFromRuntime.Mallocs))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "NextGC", float64(mFromRuntime.NextGC))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "NumForcedGC", float64(mFromRuntime.NumForcedGC))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "NumGC", float64(mFromRuntime.NumGC))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "OtherSys", float64(mFromRuntime.OtherSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "PauseTotalNs", float64(mFromRuntime.PauseTotalNs))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "StackInuse", float64(mFromRuntime.StackInuse))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "StackSys", float64(mFromRuntime.StackSys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "Sys", float64(mFromRuntime.Sys))
	sStg.GetStorage().UpdateMetricByName(constants.GaugeType, "TotalAlloc", float64(mFromRuntime.TotalAlloc))
}
