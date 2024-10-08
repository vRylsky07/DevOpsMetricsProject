// Code generated by MockGen. DO NOT EDIT.
// Source: storagetemplate.go

// Package mock_storage is a generated GoMock package.
package mock_storage

import (
	constants "DevOpsMetricsProject/internal/constants"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricsRepository is a mock of MetricsRepository interface.
type MockMetricsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsRepositoryMockRecorder
}

// MockMetricsRepositoryMockRecorder is the mock recorder for MockMetricsRepository.
type MockMetricsRepositoryMockRecorder struct {
	mock *MockMetricsRepository
}

// NewMockMetricsRepository creates a new mock instance.
func NewMockMetricsRepository(ctrl *gomock.Controller) *MockMetricsRepository {
	mock := &MockMetricsRepository{ctrl: ctrl}
	mock.recorder = &MockMetricsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsRepository) EXPECT() *MockMetricsRepositoryMockRecorder {
	return m.recorder
}

// CheckBackupStatus mocks base method.
func (m *MockMetricsRepository) CheckBackupStatus() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckBackupStatus")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckBackupStatus indicates an expected call of CheckBackupStatus.
func (mr *MockMetricsRepositoryMockRecorder) CheckBackupStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckBackupStatus", reflect.TypeOf((*MockMetricsRepository)(nil).CheckBackupStatus))
}

// GetMetricByName mocks base method.
func (m *MockMetricsRepository) GetMetricByName(mType constants.MetricType, mName string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetricByName", mType, mName)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetricByName indicates an expected call of GetMetricByName.
func (mr *MockMetricsRepositoryMockRecorder) GetMetricByName(mType, mName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetricByName", reflect.TypeOf((*MockMetricsRepository)(nil).GetMetricByName), mType, mName)
}

// IsValid mocks base method.
func (m *MockMetricsRepository) IsValid() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValid")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValid indicates an expected call of IsValid.
func (mr *MockMetricsRepositoryMockRecorder) IsValid() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValid", reflect.TypeOf((*MockMetricsRepository)(nil).IsValid))
}

// ReadMemStorageFields mocks base method.
func (m *MockMetricsRepository) ReadMemStorageFields() (map[string]float64, map[string]int) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMemStorageFields")
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(map[string]int)
	return ret0, ret1
}

// ReadMemStorageFields indicates an expected call of ReadMemStorageFields.
func (mr *MockMetricsRepositoryMockRecorder) ReadMemStorageFields() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMemStorageFields", reflect.TypeOf((*MockMetricsRepository)(nil).ReadMemStorageFields))
}

// UpdateMetricByName mocks base method.
func (m *MockMetricsRepository) UpdateMetricByName(oper constants.UpdateOperation, mType constants.MetricType, mName string, mValue float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricByName", oper, mType, mName, mValue)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricByName indicates an expected call of UpdateMetricByName.
func (mr *MockMetricsRepositoryMockRecorder) UpdateMetricByName(oper, mType, mName, mValue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricByName", reflect.TypeOf((*MockMetricsRepository)(nil).UpdateMetricByName), oper, mType, mName, mValue)
}
