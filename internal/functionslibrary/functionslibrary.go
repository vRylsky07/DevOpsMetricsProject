package functionslibrary

import (
	"DevOpsMetricsProject/internal/constants"
	"DevOpsMetricsProject/internal/coretypes"
	"DevOpsMetricsProject/internal/logger"
	"DevOpsMetricsProject/internal/storage"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
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

func EncodeBatchJSON(mStg storage.StorageInterface) (*bytes.Buffer, error) {
	gauge, counter := mStg.ReadMemStorageFields()
	length := len(gauge) + len(counter)

	mArray := make([]coretypes.Metrics, length)

	i := 0

	for k, v := range gauge {
		if i == length {
			break
		}
		mArray[i] = coretypes.Metrics{ID: k, MType: "gauge", Value: &v}
		i++
	}

	for k, v := range counter {
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
		logger.Log.Error(err.Error())
	}

	str, _ := jsonOut.ReadString(';')

	logger.Log.Info(str)

	return &jsonOut, err
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
		logger.Log.Info("gzip.NewReader failed")
		return nil, err
	}

	defer gz.Close()

	decomp, decompErr := io.ReadAll(gz)
	if decompErr != nil {
		logger.Log.Info("Decompressiong ReadAll failed")
		return nil, decompErr
	}

	newReader := strings.NewReader(string(decomp))
	newReadCloser := io.NopCloser(newReader)

	return newReadCloser, nil
}

func UpdateStorageInterfaceByMetricStruct(sStg storage.StorageInterface, mType constants.MetricType, mReceiver *coretypes.Metrics) error {
	switch mType {
	case constants.GaugeType:
		if mReceiver.Value == nil {
			return fmt.Errorf("updating gauge value pointer is nil, ID=%s", mReceiver.ID)
		}
		sStg.UpdateMetricByName(constants.RenewOperation, mType, mReceiver.ID, *mReceiver.Value)
		return nil

	case constants.CounterType:
		if mReceiver.Delta == nil {
			return fmt.Errorf("updating counter value pointer is nil, ID=%s", mReceiver.ID)
		}
		sStg.UpdateMetricByName(constants.AddOperation, mType, mReceiver.ID, float64(*mReceiver.Delta))
		return nil

	default:
		convertErr := "ConvertStringToMetricType returns NoneType"
		return errors.New(convertErr)
	}
}
