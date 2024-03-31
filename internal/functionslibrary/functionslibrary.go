package functionslibrary

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/coretypes"
	"bytes"
	"encoding/json"
	"math/rand"
	"time"
)

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

func ConvertStringToMetricType(str string) constants.MetricType {
	switch str {
	case "gauge":
		return constants.GaugeType
	case "counter":
		return constants.CounterType
	default:
		return constants.NoneType
	}
}

func ConvertMetricTypeToString(mType constants.MetricType) string {
	switch mType {
	case constants.GaugeType:
		return "gauge"
	case constants.CounterType:
		return "counter"
	default:
		return ""
	}
}

func EncodeMetricJSON(mType constants.MetricType, mName string, mValue float64) (*bytes.Buffer, error) {
	src := coretypes.Metrics{}

	src.ID = mName
	src.MType = ConvertMetricTypeToString(mType)

	switch mType {
	case constants.GaugeType:
		src.Value = &mValue
	case constants.CounterType:
		intValue := int64(mValue)
		src.Delta = &intValue
	}

	var jsonOut bytes.Buffer

	jsonEncoder := json.NewEncoder(&jsonOut)
	err := jsonEncoder.Encode(src)

	return &jsonOut, err
}
