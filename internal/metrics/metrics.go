package metrics

import (
	"DevOpsMetricsProject/internal/functionslibrary"
	"runtime"
)

type MetricsCollectorInterface interface {
	InitMetricsCollector() MetricsCollectorInterface
	UpdateCounterMetrics()
	UpdateGaugeMetrics()
	ReadMetricsCollector() (map[string]float64, map[string]int)
}

type MetricsCollector struct {
	gauge   map[string]float64
	counter map[string]int
}

func (mCollector *MetricsCollector) InitMetricsCollector() MetricsCollectorInterface {
	mCollector.gauge, mCollector.counter = map[string]float64{}, map[string]int{}
	mCollector.counter["PollCount"] = 0
	return mCollector
}

// сбор всех counter-метрик
func (mCollector *MetricsCollector) UpdateCounterMetrics() {
	if mCollector == nil {
		mCollector.InitMetricsCollector()
	}

	mCollector.counter["PollCount"] += 1
}

// сбор всех gauge-метрик
func (mCollector *MetricsCollector) UpdateGaugeMetrics() {

	if mCollector == nil {
		mCollector.InitMetricsCollector()
	}

	var mFromRuntime *runtime.MemStats = &runtime.MemStats{}
	runtime.ReadMemStats(mFromRuntime)

	mCollector.gauge["RandomValue"] = functionslibrary.GenerateRandomValue(-10000, 10000, 3)

	mCollector.gauge["Alloc"] = float64(mFromRuntime.Alloc)
	mCollector.gauge["BuckHashSys"] = float64(mFromRuntime.BuckHashSys)
	mCollector.gauge["Frees"] = float64(mFromRuntime.Frees)
	mCollector.gauge["GCCPUFraction"] = float64(mFromRuntime.GCCPUFraction)
	mCollector.gauge["GCSys"] = float64(mFromRuntime.GCSys)
	mCollector.gauge["HeapAlloc"] = float64(mFromRuntime.HeapAlloc)
	mCollector.gauge["HeapIdle"] = float64(mFromRuntime.HeapIdle)
	mCollector.gauge["HeapInuse"] = float64(mFromRuntime.HeapInuse)
	mCollector.gauge["HeapReleased"] = float64(mFromRuntime.HeapReleased)
	mCollector.gauge["HeapObjects"] = float64(mFromRuntime.HeapObjects)
	mCollector.gauge["HeapSys"] = float64(mFromRuntime.HeapSys)
	mCollector.gauge["LastGC"] = float64(mFromRuntime.LastGC)
	mCollector.gauge["Lookups"] = float64(mFromRuntime.Lookups)
	mCollector.gauge["MCacheInuse"] = float64(mFromRuntime.MCacheInuse)
	mCollector.gauge["MCacheSys"] = float64(mFromRuntime.MCacheSys)
	mCollector.gauge["MSpanInuse"] = float64(mFromRuntime.MSpanInuse)
	mCollector.gauge["MSpanSys"] = float64(mFromRuntime.MSpanSys)
	mCollector.gauge["Mallocs"] = float64(mFromRuntime.Mallocs)
	mCollector.gauge["NextGC"] = float64(mFromRuntime.NextGC)
	mCollector.gauge["NumForcedGC"] = float64(mFromRuntime.NumForcedGC)
	mCollector.gauge["NumGC"] = float64(mFromRuntime.NumGC)
	mCollector.gauge["OtherSys"] = float64(mFromRuntime.OtherSys)
	mCollector.gauge["PauseTotalNs"] = float64(mFromRuntime.PauseTotalNs)
	mCollector.gauge["StackInuse"] = float64(mFromRuntime.StackInuse)
	mCollector.gauge["StackSys"] = float64(mFromRuntime.StackSys)
	mCollector.gauge["Sys"] = float64(mFromRuntime.Sys)
	mCollector.gauge["TotalAlloc"] = float64(mFromRuntime.TotalAlloc)
}

func (mCollector *MetricsCollector) ReadMetricsCollector() (map[string]float64, map[string]int) {
	return mCollector.gauge, mCollector.counter
}
