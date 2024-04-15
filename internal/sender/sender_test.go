package sender

import (
	"DevOpsMetricsProject/internal/constants"
	storage_custom_mocks "DevOpsMetricsProject/internal/storage/custom_mock"
	mock_storage "DevOpsMetricsProject/internal/storage/mocks"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSenderStorage_InitSenderStorage(t *testing.T) {
	tests := []struct {
		name   string
		actual *SenderStorage
		args   *mock_storage.MockStorageInterface
	}{
		{
			name:   "Sender storage initialization",
			actual: &SenderStorage{},
			args:   &mock_storage.MockStorageInterface{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.InitSenderStorage(nil, tt.args)
			assert.NotNil(t, tt.actual)
			assert.NotNil(t, tt.actual.GetStorage())
		})
	}
}

func TestSenderStorage_updateGaugeMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockStorageInterface(mockCtrl)

	tests := []struct {
		name   string
		actual *SenderStorage
	}{{
		name:   "updateGaugeMetrics",
		actual: &SenderStorage{storageChecker, false, nil},
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
		actual *SenderStorage
	}{{
		name:   "updateCounterMetrics",
		actual: &SenderStorage{storageChecker, false, nil},
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
		actual *SenderStorage
	}{{
		name:   "Update Metrics",
		actual: &SenderStorage{storageChecker, false, nil},
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.GaugeType, gomock.Any(), gomock.Any()).AnyTimes()
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.CounterType, gomock.Any(), gomock.Any()).AnyTimes().Return()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		sStg *SenderStorage
	}{{
		name: "Send metrics to server HTTP",
		sStg: &SenderStorage{testingStorage, false, nil},
	},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
