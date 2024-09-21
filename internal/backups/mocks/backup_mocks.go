// Code generated by MockGen. DO NOT EDIT.
// Source: backup.go

// Package mock_backup is a generated GoMock package.
package mock_backup

import (
	constants "DevOpsMetricsProject/internal/constants"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricsBackup is a mock of MetricsBackup interface.
type MockMetricsBackup struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsBackupMockRecorder
}

// MockMetricsBackupMockRecorder is the mock recorder for MockMetricsBackup.
type MockMetricsBackupMockRecorder struct {
	mock *MockMetricsBackup
}

// NewMockMetricsBackup creates a new mock instance.
func NewMockMetricsBackup(ctrl *gomock.Controller) *MockMetricsBackup {
	mock := &MockMetricsBackup{ctrl: ctrl}
	mock.recorder = &MockMetricsBackupMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsBackup) EXPECT() *MockMetricsBackupMockRecorder {
	return m.recorder
}

// GetAllData mocks base method.
func (m *MockMetricsBackup) GetAllData() (*map[string]float64, *map[string]int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllData")
	ret0, _ := ret[0].(*map[string]float64)
	ret1, _ := ret[1].(*map[string]int)
	return ret0, ret1
}

// GetAllData indicates an expected call of GetAllData.
func (mr *MockMetricsBackupMockRecorder) GetAllData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllData", reflect.TypeOf((*MockMetricsBackup)(nil).GetAllData))
}

// IsValid mocks base method.
func (m *MockMetricsBackup) IsValid() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValid")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValid indicates an expected call of IsValid.
func (mr *MockMetricsBackupMockRecorder) IsValid() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValid", reflect.TypeOf((*MockMetricsBackup)(nil).IsValid))
}

// UpdateMetricBackup mocks base method.
func (m *MockMetricsBackup) UpdateMetricBackup(mType constants.MetricType, mName string, mValue float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricBackup", mType, mName, mValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricBackup indicates an expected call of UpdateMetricBackup.
func (mr *MockMetricsBackupMockRecorder) UpdateMetricBackup(mType, mName, mValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricBackup", reflect.TypeOf((*MockMetricsBackup)(nil).UpdateMetricBackup), mType, mName, mValue)
}

// MockPingerDB is a mock of PingerDB interface.
type MockPingerDB struct {
	ctrl     *gomock.Controller
	recorder *MockPingerDBMockRecorder
}

// MockPingerDBMockRecorder is the mock recorder for MockPingerDB.
type MockPingerDBMockRecorder struct {
	mock *MockPingerDB
}

// NewMockPingerDB creates a new mock instance.
func NewMockPingerDB(ctrl *gomock.Controller) *MockPingerDB {
	mock := &MockPingerDB{ctrl: ctrl}
	mock.recorder = &MockPingerDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPingerDB) EXPECT() *MockPingerDBMockRecorder {
	return m.recorder
}

// PingDB mocks base method.
func (m *MockPingerDB) PingDB() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingDB")
	ret0, _ := ret[0].(error)
	return ret0
}

// PingDB indicates an expected call of PingDB.
func (mr *MockPingerDBMockRecorder) PingDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingDB", reflect.TypeOf((*MockPingerDB)(nil).PingDB))
}