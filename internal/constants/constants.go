package constants

type MetricType int

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

type SaveMode int

const (
	DatabaseMode SaveMode = iota
	FileMode
	InMemoryMode
)
