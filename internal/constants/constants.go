package constants

// Enum для типа метрики
type MetricType int

// константы для проверок типа
const (
	GaugeType MetricType = iota
	CounterType
	NoneType
)

type UpdateOperation int

const (
	AddOperation UpdateOperation = iota
	RenewOperation
)

type DecimalCount int

const (
	CounterDecimal  = 0
	GaugeDecimal    = 3
	NoneTypeDecimal = 0
)
