package metrics

import (
	"DevOpsMetricsProject/internal/storage"
	"math/rand"
	"runtime"
	"time"
)

// собираем все метрики для отправки
func GetAllMetrics() storage.MemStorage {
	MemStg := storage.MemStorage{}
	MemStg.SetMemStorage(GetGaugeMetrics(), GetCounterMetrics())
	return MemStg
}

// создание экземпляра счетчика
var PollCount int = 0

// обновление счётчика
func UpdatePollCounter() int {
	PollCount++
	return PollCount
}

// сбор всех counter-метрик
func GetCounterMetrics() map[string]int {
	finalCounterMap := map[string]int{}
	finalCounterMap["PollCount"] += UpdatePollCounter()
	return finalCounterMap
}

// сбор всех gauge-метрик
func GetGaugeMetrics() map[string]float64 {
	finalGaugeMap := map[string]float64{}

	finalGaugeMap["RandomValue"] = GenerateRandomValue(-10000, 10000, 3)

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
func GenerateRandomValue(min int, max int, precision int) float64 {
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
