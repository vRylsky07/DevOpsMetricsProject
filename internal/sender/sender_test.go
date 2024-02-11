package sender

import (
	"DevOpsMetricsProject/internal/constants"
	mock_storage "DevOpsMetricsProject/internal/storage/mocks"
	"testing"
	"time"

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
			tt.actual.InitSenderStorage(tt.args)
			assert.NotNil(t, tt.actual)
			assert.NotNil(t, tt.actual.GetStorage())
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
		actual: &SenderStorage{storageChecker, false},
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(constants.GaugeType, gomock.Any(), gomock.Any()).AnyTimes()
	storageChecker.EXPECT().UpdateMetricByName(constants.CounterType, gomock.Any(), gomock.Any()).AnyTimes()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go tt.actual.UpdateMetrics(0)
			time.Sleep(time.Duration(1 * time.Second))
			tt.actual.StopAgentProcessing()
			assert.NotNil(t, tt.actual)

			//assert.NotNil(t, c)
		})
	}
}
