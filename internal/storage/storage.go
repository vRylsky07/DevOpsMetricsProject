package storage

// Enum для типа метрики
type MetricType int

// константы для проверок типа
const (
	Gauge MetricType = iota
	Counter
)

// описание хранилища данных
type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int
}

func (mStg *MemStorage) SetMemStorage(g map[string]float64, c map[string]int) {
	mStg.gauge = g
	mStg.counter = c
}

// инстанциация хранилища
var mStrg *MemStorage = &MemStorage{
	gauge:   map[string]float64{},
	counter: map[string]int{},
}

// обновление метрик с указанием Enum типа
func (mStr *MemStorage) UpdateMetricByName(mType MetricType, mName string, mValue float64) {
	switch mType {
	case Gauge:
		mStr.gauge[mName] = mValue
	case Counter:
		mStr.counter[mName] += int(mValue)
	}
}

// геттер метрик по имени
func (mStr *MemStorage) GetMetricByName(mType MetricType, mName string) float64 {
	switch mType {
	case Gauge:
		gMetric, wasFound := mStr.gauge[mName]
		if !wasFound {
			return 0.0
		}
		return gMetric
	case Counter:
		cMetric, wasFound := mStr.counter[mName]
		if !wasFound {
			return 0.0
		}
		return float64(cMetric)
	default:
		return 0.0
	}
}

// геттер экземпляра хранилища
func GetMemStorage() *MemStorage {
	return mStrg
}
