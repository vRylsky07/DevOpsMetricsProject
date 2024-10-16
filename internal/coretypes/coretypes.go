package coretypes

import "bytes"

//go:generate easyjson -all coretypes.go

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type ReqProps struct {
	URL        string
	Body       *bytes.Buffer
	Sign       string
	MetricName string
	IsBatch    bool
}
