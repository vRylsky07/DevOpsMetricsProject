package sender

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/metrics"
	"DevOpsMetricsProject/internal/storage"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type SenderInterface interface {
	InitSenderStorage(newStg storage.StorageInterface)
	GetStorage() storage.StorageInterface
	UpdateMetrics(mCollector metrics.MetricsCollectorInterface, pollInterval int)
	SendMetricsHTTP(reportInterval int)
	CreateMetricURL(mType constants.MetricType, mainURL string, name string, value float64) string
}

type SenderStorage struct {
	senderMemStorage storage.StorageInterface
}

func (sStg *SenderStorage) GetStorage() storage.StorageInterface {
	return sStg.senderMemStorage
}

func (sStg *SenderStorage) InitSenderStorage(newStg storage.StorageInterface) {
	sStg.senderMemStorage = newStg
}

// обновление метрик
func (sStg *SenderStorage) UpdateMetrics(mCollector metrics.MetricsCollectorInterface, pollInterval int) {

	if sStg == nil {
		sStg.GetStorage().SetMemStorage(map[string]float64{}, map[string]int{})
	}

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		mCollector.UpdateCounterMetrics()
		mCollector.UpdateGaugeMetrics()
		gauge, counter := mCollector.ReadMetricsCollector()
		sStg.GetStorage().SetMemStorage(gauge, counter)
	}
}

// отправка метрик
func (sStg *SenderStorage) SendMetricsHTTP(reportInterval int) {

	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		if sStg == nil {
			return
		}

		gauge, counter := sStg.GetStorage().ReadMemStorageFields()

		for nameGauge, valueGauge := range gauge {
			finalURL := sStg.CreateMetricURL(constants.GaugeType, "http://localhost:8080", nameGauge, valueGauge)
			resp, err := http.Post(sStg.CreateMetricURL(constants.GaugeType, finalURL, nameGauge, valueGauge), "text/plain", nil)
			if err != nil {
				fmt.Println("Server is not responding. URL to send was: " + finalURL)
				continue
			}
			fmt.Println(finalURL, "Gauge update was successful! Status code: ", resp.StatusCode)
		}

		for nameCounter, valueCounter := range counter {
			finalURL := sStg.CreateMetricURL(constants.CounterType, "http://localhost:8080", nameCounter, float64(valueCounter))
			resp, err := http.Post(finalURL, "text/plain", nil)
			if err != nil {
				fmt.Println("Server is not responding. URL to send was: " + finalURL)
				continue
			}
			fmt.Println(finalURL, " Counter update was successful! Status code: ", resp.StatusCode)
		}
	}
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
	return mainURL + "/update" + mTypeString + name + "/" + strconv.FormatFloat(value, 'f', 6, 64)
}

func CreateSender() *SenderStorage {
	agentStorage := storage.MemStorage{}
	agentStorage.InitMemStorage()

	mSender := &SenderStorage{}
	mSender.InitSenderStorage(&agentStorage)

	return mSender
}
