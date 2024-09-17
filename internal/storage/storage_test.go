package storage

/*
import (
	"DevOpsMetricsProject/internal/constants"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateMetricByName(t *testing.T) {
	type args struct {
		oper   constants.UpdateOperation
		mType  constants.MetricType
		mName  string
		mValue float64
	}
	tests := []struct {
		name string
		mStg *MemStorage
		args *args
		want float64
	}{
		{
			name: "UpdateMetricByName: Renew, Gauge case",
			mStg: &MemStorage{},
			args: &args{constants.RenewOperation, constants.GaugeType, "TestGauge", 56.662},
			want: 56.662,
		},
		{
			name: "UpdateMetricByName: Add, Gauge case",
			mStg: &MemStorage{},
			args: &args{constants.AddOperation, constants.GaugeType, "TestGauge", 70.5},
			want: 100.5,
		},
		{
			name: "UpdateMetricByName: Renew, Counter case",
			mStg: &MemStorage{},
			args: &args{constants.RenewOperation, constants.CounterType, "TestGauge", 70},
			want: 70,
		},
		{
			name: "UpdateMetricByName: Add, Counter case",
			mStg: &MemStorage{},
			args: &args{constants.AddOperation, constants.CounterType, "TestGauge", 70},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mStg.InitMemStorage()
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

func TestMemStorage_GetMetricByName(t *testing.T) {
	type args struct {
		mType constants.MetricType
		mName string
	}
	tests := []struct {
		name    string
		mStg    *MemStorage
		args    *args
		want    float64
		wantErr bool
	}{
		{
			name:    "GetMetricByName: Gauge, Name is valid",
			mStg:    &MemStorage{},
			args:    &args{constants.GaugeType, "TestGauge"},
			want:    77.53,
			wantErr: false,
		},
		{
			name:    "GetMetricByName: Gauge, Name is not valid",
			mStg:    &MemStorage{},
			args:    &args{constants.GaugeType, "TestGaugeInvalid"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "GetMetricByName: Counter, Name is valid",
			mStg:    &MemStorage{},
			args:    &args{constants.CounterType, "TestCounter"},
			want:    31,
			wantErr: false,
		},
		{
			name:    "GetMetricByName: Counter, Name is not valid",
			mStg:    &MemStorage{},
			args:    &args{constants.CounterType, "TestCounterInvalid"},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mStg.InitMemStorage()
			tt.mStg.counter["TestCounter"] = 31
			tt.mStg.gauge["TestGauge"] = 77.53

			gotErr := false
			got, err := tt.mStg.GetMetricByName(tt.args.mType, tt.args.mName)
			if err != nil {
				gotErr = true
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
*/
