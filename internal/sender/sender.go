package sender

import (
	"DevOpsMetricsProject/internal/metrics"
	"DevOpsMetricsProject/internal/storage"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var MemStg *storage.MemStorage = storage.CreateMemStorage()

// обновление метрик
func UpdateMetrics(pollInterval int) *storage.MemStorage {

	if MemStg == nil {
		MemStg.SetMemStorage(map[string]float64{}, map[string]int{})
	}

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		MemStg.SetMemStorage(metrics.GetGaugeMetrics(), metrics.GetCounterMetrics())
	}
}

// отправка метрик
func SendMetricsHTTP(reportInterval int) {
	for {
		time.Sleep(time.Duration(reportInterval) * time.Second)
		if MemStg == nil {
			return
		}

		gauge, counter := MemStg.ReadMemStorageFields()

		for nameGauge, valueGauge := range gauge {
			finalURL := CreateMetricURL(storage.Gauge, "http://localhost:8080", nameGauge, valueGauge)
			resp, err := http.Post(CreateMetricURL(storage.Gauge, finalURL, nameGauge, valueGauge), "text/plain", nil)
			if err != nil {
				fmt.Println("Server is not responding. URL to send was: " + finalURL)
				continue
			}
			fmt.Println(finalURL, "Gauge update was successful! Status code: ", resp.StatusCode)
		}

		for nameCounter, valueCounter := range counter {
			finalURL := CreateMetricURL(storage.Gauge, "http://localhost:8080", nameCounter, float64(valueCounter))
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
func CreateMetricURL(mType storage.MetricType, mainURL string, name string, value float64) string {
	mTypeString := ""

	switch mType {
	case storage.Gauge:
		mTypeString = "/gauge/"
	case storage.Counter:
		mTypeString = "/counter/"
	}
	return mainURL + "/update" + mTypeString + name + "/" + strconv.FormatFloat(value, 'f', 6, 64)
}
