package sender

import (
	"DevOpsMetricsProject/internal/metrics"
	"DevOpsMetricsProject/internal/storage"
	"time"
)

var MemStg *storage.MemStorage = storage.CreateMemStorage()

func UpdateMetrics(pollInterval int) *storage.MemStorage {

	if MemStg == nil {
		MemStg.SetMemStorage(map[string]float64{}, map[string]int{})
	}

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		MemStg.SetMemStorage(metrics.GetGaugeMetrics(), metrics.GetCounterMetrics())
	}
}

func SendMetricsHTTP(reportInterval int) {

}
