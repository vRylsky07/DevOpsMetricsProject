package constants

// Enum для типа метрики
type MetricType int

// константы для проверок типа
const (
	GaugeType MetricType = iota
	CounterType
)
