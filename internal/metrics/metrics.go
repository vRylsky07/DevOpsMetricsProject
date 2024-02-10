package metrics

import (
	"math/rand"
	"runtime"
	"time"
)

type MetricsCollectorInterface interface {
	InitMetricsCollector() MetricsCollectorInterface
	UpdateCounterMetrics() map[string]int
	UpdateGaugeMetrics() map[string]float64
	GenerateRandomValue(min int, max int, precision int) float64
}

type MetricsCollector struct {
	PollCount int
}

func (mCollector *MetricsCollector) InitMetricsCollector() MetricsCollectorInterface {
	mCollector.PollCount = 0
	return mCollector
}

// сбор всех counter-метрик
func (mCollector *MetricsCollector) UpdateCounterMetrics() map[string]int {
	finalCounterMap := map[string]int{}
	mCollector.PollCount++
	finalCounterMap["PollCount"] += mCollector.PollCount
	return finalCounterMap
}

// сбор всех gauge-метрик
func (mCollector *MetricsCollector) UpdateGaugeMetrics() map[string]float64 {
	finalGaugeMap := map[string]float64{}

	finalGaugeMap["RandomValue"] = mCollector.GenerateRandomValue(-10000, 10000, 3)

	var mFromRuntime *runtime.MemStats = &runtime.MemStats{}
	runtime.ReadMemStats(mFromRuntime)

	finalGaugeMap["Alloc"] = float64(mFromRuntime.Alloc)
	finalGaugeMap["BuckHashSys"] = float64(mFromRuntime.BuckHashSys)
	finalGaugeMap["Frees"] = float64(mFromRuntime.Frees)
	finalGaugeMap["GCCPUFraction"] = float64(mFromRuntime.GCCPUFraction)
	finalGaugeMap["GCSys"] = float64(mFromRuntime.GCSys)
	finalGaugeMap["HeapAlloc"] = float64(mFromRuntime.HeapAlloc)
	finalGaugeMap["HeapIdle"] = float64(mFromRuntime.HeapIdle)
	finalGaugeMap["HeapInuse"] = float64(mFromRuntime.HeapInuse)
	finalGaugeMap["HeapReleased"] = float64(mFromRuntime.HeapReleased)
	finalGaugeMap["HeapObjects"] = float64(mFromRuntime.HeapObjects)
	finalGaugeMap["HeapSys"] = float64(mFromRuntime.HeapSys)
	finalGaugeMap["LastGC"] = float64(mFromRuntime.LastGC)
	finalGaugeMap["Lookups"] = float64(mFromRuntime.Lookups)
	finalGaugeMap["MCacheInuse"] = float64(mFromRuntime.MCacheInuse)
	finalGaugeMap["MCacheSys"] = float64(mFromRuntime.MCacheSys)
	finalGaugeMap["MSpanInuse"] = float64(mFromRuntime.MSpanInuse)
	finalGaugeMap["MSpanSys"] = float64(mFromRuntime.MSpanSys)
	finalGaugeMap["Mallocs"] = float64(mFromRuntime.Mallocs)
	finalGaugeMap["NextGC"] = float64(mFromRuntime.NextGC)
	finalGaugeMap["NumForcedGC"] = float64(mFromRuntime.NumForcedGC)
	finalGaugeMap["NumGC"] = float64(mFromRuntime.NumGC)
	finalGaugeMap["OtherSys"] = float64(mFromRuntime.OtherSys)
	finalGaugeMap["PauseTotalNs"] = float64(mFromRuntime.PauseTotalNs)
	finalGaugeMap["StackInuse"] = float64(mFromRuntime.StackInuse)
	finalGaugeMap["StackSys"] = float64(mFromRuntime.StackSys)
	finalGaugeMap["Sys"] = float64(mFromRuntime.Sys)
	finalGaugeMap["TotalAlloc"] = float64(mFromRuntime.TotalAlloc)

	return finalGaugeMap
}

// генерация рандомного значения
func (mCollector *MetricsCollector) GenerateRandomValue(min int, max int, precision int) float64 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	intPart := rng.Intn(max-min+1) + min

	decimalPart := float64(int(rand.Float64()*float64((precision*10)))) / float64((precision * 10))
	mixedValue := float64(intPart) + decimalPart

	if mixedValue <= float64(min) || mixedValue >= float64(max) {
		return float64(intPart)
	}

	return mixedValue
}
