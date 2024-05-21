package sender

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	storage_custom_mocks "DevOpsMetricsProject/internal/storage/custom_mock"
	mock_storage "DevOpsMetricsProject/internal/storage/mocks"

	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSenderStorage_updateGaugeMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockStorageInterface(mockCtrl)

	tests := []struct {
		name   string
		actual *dompsender
	}{{
		name:   "updateGaugeMetrics",
		actual: &dompsender{storageChecker, false, &configs.ClientConfig{}},
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.GaugeType, gomock.Any(), gomock.Any()).Times(28)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.updateGaugeMetrics()
		})
	}
}

func TestSenderStorage_updateCounterMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockStorageInterface(mockCtrl)

	tests := []struct {
		name   string
		actual *dompsender
	}{{
		name:   "updateCounterMetrics",
		actual: &dompsender{storageChecker, false, &configs.ClientConfig{}},
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.CounterType, gomock.Any(), gomock.Any()).Times(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.updateCounterMetrics()
		})
	}
}

func TestSenderStorage_UpdateMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockStorageInterface(mockCtrl)

	tests := []struct {
		name   string
		actual *dompsender
	}{{
		name:   "Update Metrics",
		actual: &dompsender{storageChecker, false, &configs.ClientConfig{}},
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.GaugeType, gomock.Any(), gomock.Any()).AnyTimes()
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.CounterType, gomock.Any(), gomock.Any()).AnyTimes().Return()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.cfg.PollInterval = -1
			tt.actual.UpdateMetrics()
		})
	}
}

func TestSenderStorage_SendMetricsHTTP(t *testing.T) {
	g := map[string]float64{"testGauge": 1}
	c := map[string]int{"testCounter": 2}

	testingStorage := &storage_custom_mocks.StorageMockCustom{Gauge: g, Counter: c}

	tests := []struct {
		name string
		sStg *dompsender
	}{{
		name: "Send metrics to server HTTP",
		sStg: &dompsender{testingStorage, false, &configs.ClientConfig{}},
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sStg.cfg.ReportInterval = -1
			errs := tt.sStg.SendMetricsHTTP()

			g, c := tt.sStg.GetStorage().ReadMemStorageFields()
			if len(errs) != (len(g) + len(c)) {
				t.Errorf("Send metrics to server HTTP FAILED!")
			}

			var isWantedErr bool
			for _, err := range errs {
				isWantedErr = strings.Contains(fmt.Sprint(err), "Server is not responding")
				if isWantedErr == false {
					break
				}
			}

			assert.True(t, isWantedErr)
		})
	}
}
