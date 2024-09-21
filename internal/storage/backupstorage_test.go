package storage

import (
	backup_mocks "DevOpsMetricsProject/internal/backups/mocks"
	"DevOpsMetricsProject/internal/constants"
	logger_mocks "DevOpsMetricsProject/internal/logger/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBackupSupportStorage_UpdateMetricByName(t *testing.T) {
	MockCtrl := gomock.NewController(t)
	backuper := backup_mocks.NewMockMetricsBackup(MockCtrl)

	LogsMockCtrl := gomock.NewController(t)
	log := logger_mocks.NewMockRecorder(LogsMockCtrl)

	type args struct {
		oper   constants.UpdateOperation
		mType  constants.MetricType
		mName  string
		mValue float64
	}
	tests := []struct {
		name string
		mStg *BackupSupportStorage
		args *args
		want float64
	}{
		{
			name: "UpdateMetricByName: Renew, Gauge case",
			mStg: nil,
			args: &args{constants.RenewOperation, constants.GaugeType, "TestGauge", 56.662},
			want: 56.662,
		},
		{
			name: "UpdateMetricByName: Add, Gauge case",
			mStg: nil,
			args: &args{constants.AddOperation, constants.GaugeType, "TestGauge", 70.5},
			want: 100.5,
		},
		{
			name: "UpdateMetricByName: Renew, Counter case",
			mStg: nil,
			args: &args{constants.RenewOperation, constants.CounterType, "TestGauge", 70},
			want: 70,
		},
		{
			name: "UpdateMetricByName: Add, Counter case",
			mStg: nil,
			args: &args{constants.AddOperation, constants.CounterType, "TestGauge", 70},
			want: 100,
		},
	}

	backuper.EXPECT().UpdateMetricBackup(gomock.Any(), gomock.Any(), gomock.Any()).Times(len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mStg := &BackupSupportStorage{}
			mStg.backup = backuper
			mStg.log = log
			mStg.gauge, mStg.counter = map[string]float64{}, map[string]int{}

			tt.mStg = mStg
			tt.mStg.counter[tt.args.mName] = 30
			tt.mStg.gauge[tt.args.mName] = 30
			tt.mStg.UpdateMetricByName(tt.args.oper, tt.args.mType, tt.args.mName, tt.args.mValue)

			updatedValue, isExist := tt.mStg.gauge[tt.args.mName]

			if tt.args.mType == constants.CounterType {
				var updatedCounter int
				updatedCounter, isExist = tt.mStg.counter[tt.args.mName]
				updatedValue = float64(updatedCounter)
			}

			if isExist {
				assert.Equal(t, tt.want, updatedValue)
			}
			assert.True(t, isExist)
		})
	}
}
