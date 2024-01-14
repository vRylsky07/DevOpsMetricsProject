package storage

type MetricType int

const (
	Gauge MetricType = iota
	Counter
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int
}

func (mStr *MemStorage) UpdateMetricByName(mType MetricType, mName string, mValue float64) {
	switch mType {
	case Gauge:
		mStr.gauge[mName] = mValue
	case Counter:
		mStr.counter[mName] += int(mValue)
	}
}

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

func InitMemStorage() {
	mStrg = &MemStorage{
		gauge:   map[string]float64{},
		counter: map[string]int{},
	}
}

func GetMemStrorage() *MemStorage {
	return mStrg
}

var mStrg *MemStorage
