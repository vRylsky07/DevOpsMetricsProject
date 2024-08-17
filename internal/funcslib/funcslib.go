package funcslib

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/coretypes"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"strings"
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

func DecodeMetricJSON(req io.ReadCloser) (*coretypes.Metrics, error) {
	var mReceiver coretypes.Metrics
	err := json.NewDecoder(req).Decode(&mReceiver)
	return &mReceiver, err
}

func EncodeBatchJSON(gauge *map[string]float64, counter *map[string]int) (*bytes.Buffer, error) {

	if gauge == nil || counter == nil {
		return nil, errors.New("EncodeBatchJSON() failed. One of two maps is empty")
	}

	length := len(*gauge) + len(*counter)

	mArray := make([]coretypes.Metrics, length)

	i := 0

	for k, v := range *gauge {
		if i == length {
			break
		}
		storeValue := v
		mArray[i] = coretypes.Metrics{ID: k, MType: "gauge", Value: &storeValue}
		i++
	}

	for k, v := range *counter {
		if i == length {
			break
		}
		delta := int64(v)
		mArray[i] = coretypes.Metrics{ID: k, MType: "counter", Delta: &delta}
		i++
	}

	var jsonOut bytes.Buffer
	jsonEncoder := json.NewEncoder(&jsonOut)
	err := jsonEncoder.Encode(mArray)

	if err != nil {
		return nil, err
	}

	return &jsonOut, err
}

func DecodeBatchJSON(req io.ReadCloser) (*[]coretypes.Metrics, error) {
	var mReceiver []coretypes.Metrics
	err := json.NewDecoder(req).Decode(&mReceiver)
	return &mReceiver, err
}

func CompressData(data []byte) (*bytes.Buffer, error) {
	var b bytes.Buffer

	c, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)

	if err != nil {
		return nil, err
	}

	_, err = c.Write(data)

	if err != nil {
		return nil, err
	}

	err = c.Close()

	if err != nil {
		return nil, err
	}

	return &b, nil
}

func DecompressData(body io.ReadCloser) (io.ReadCloser, error) {
	gz, err := gzip.NewReader(body)
	if err != nil {
		return nil, err
	}

	defer gz.Close()

	decomp, decompErr := io.ReadAll(gz)
	if decompErr != nil {
		return nil, decompErr
	}

	newReader := strings.NewReader(string(decomp))
	newReadCloser := io.NopCloser(newReader)

	return newReadCloser, nil
}
