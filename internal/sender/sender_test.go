package sender

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	mock_storage "DevOpsMetricsProject/internal/storage/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/stretchr/testify/assert"
)

func Test_dompsender_updateGaugeMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockMetricsRepository(mockCtrl)
	sender, err := CreateSender(&configs.ClientConfig{}, 10)
	assert.Nil(t, err)
	sender.senderMemStorage = storageChecker

	tests := []struct {
		name   string
		actual *dompsender
	}{{
		name:   "updateGaugeMetrics",
		actual: sender,
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.GaugeType, gomock.Any(), gomock.Any()).Times(28)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.updateGaugeMetrics()
		})
	}
}

func Test_dompsender_updateCounterMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockMetricsRepository(mockCtrl)
	sender, err := CreateSender(&configs.ClientConfig{}, 10)
	assert.Nil(t, err)
	sender.senderMemStorage = storageChecker

	tests := []struct {
		name   string
		actual *dompsender
	}{{
		name:   "updateGaugeMetrics",
		actual: sender,
	},
	}
	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), constants.CounterType, gomock.Any(), gomock.Any()).Times(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.updateCounterMetrics()
		})
	}
}

func Test_dompsender_UpdateMetrics(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	storageChecker := mock_storage.NewMockMetricsRepository(mockCtrl)
	sender, err := CreateSender(&configs.ClientConfig{}, 10)
	assert.Nil(t, err)
	sender.senderMemStorage = storageChecker

	tests := []struct {
		name   string
		actual *dompsender
	}{{
		name:   "updateGaugeMetrics",
		actual: sender,
	},
	}

	count, err := cpu.Counts(true)

	if err != nil {
		t.Error(err.Error())
	}

	storageChecker.EXPECT().UpdateMetricByName(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(31 + count)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.actual.cfg.PollInterval = -1
			tt.actual.UpdateMetrics()
		})
	}
}

func Test_dompsender_SendMetricsHTTP(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	testingStorage := mock_storage.NewMockMetricsRepository(mockCtrl)
	sender, err := CreateSender(&configs.ClientConfig{}, 10)
	assert.Nil(t, err)
	sender.senderMemStorage = testingStorage

	tests := []struct {
		name string
		sStg *dompsender
	}{{
		name: "Send metrics to server HTTP",
		sStg: sender,
	},
	}

	testingStorage.EXPECT().ReadMemStorageFields().Times(1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sStg.cfg.ReportInterval = -1
			tt.sStg.SendMetricsHTTP()
		})
	}
}
