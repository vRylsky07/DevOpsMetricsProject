package storage

import (
	"reflect"
	"testing"
)

func TestCreateMemStorage(t *testing.T) {
	tests := []struct {
		name string
		want *MemStorage
	}{
		{
			name: "CreateMemStorageSuccess",
			want: &MemStorage{map[string]float64{}, map[string]int{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateMemStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMemStorage failed! Got: %v, Wanted:  %v", got, tt.want)
			}
		})
	}
}

// go test DevOpsMetricsProject\internal\storage\ -count 1 -v

func TestMemStorage_SetMemStorage(t *testing.T) {

	tests := []struct {
		name string
		mStg *MemStorage
		args *MemStorage
		want *MemStorage
	}{
		{
			name: "SetMemStorageSuccess",
			mStg: &MemStorage{map[string]float64{}, map[string]int{}},
			args: &MemStorage{map[string]float64{"testGauge": 64.5}, map[string]int{"testCounter": 34}},
			want: &MemStorage{map[string]float64{"testGauge": 64.5}, map[string]int{"testCounter": 34}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mStg.SetMemStorage(tt.args.gauge, tt.args.counter)
			if !reflect.DeepEqual(tt.mStg, tt.want) {
				t.Errorf("SetMemStorage failed! Got: %v, Wanted:  %v", tt.mStg, tt.want)
			}
		})
	}
}

func TestMemStorage_ReadMemStorageFields(t *testing.T) {
	tests := []struct {
		name  string
		mStg  *MemStorage
		wantG map[string]float64
		wantC map[string]int
	}{
		{
			name:  "ReadMemStorageFieldsSuccess",
			mStg:  &MemStorage{map[string]float64{"testGauge": 64.5}, map[string]int{"testCounter": 34}},
			wantG: map[string]float64{"testGauge": 64.5},
			wantC: map[string]int{"testCounter": 34},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotG, gotC := tt.mStg.ReadMemStorageFields()
			if !reflect.DeepEqual(gotG, tt.wantG) {
				t.Errorf("MemStorage.ReadMemStorageFields() gotG = %v, want %v", gotG, tt.wantG)
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("MemStorage.ReadMemStorageFields() gotC = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}
