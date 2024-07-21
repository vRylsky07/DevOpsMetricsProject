// Code generated by MockGen. DO NOT EDIT.
// Source: sender.go

// Package mock_sender is a generated GoMock package.
package mock_sender

import (
	constants "DevOpsMetricsProject/internal/constants"
	storage "DevOpsMetricsProject/internal/storage"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMetricsProvider is a mock of MetricsProvider interface.
type MockMetricsProvider struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsProviderMockRecorder
}

// MockMetricsProviderMockRecorder is the mock recorder for MockMetricsProvider.
type MockMetricsProviderMockRecorder struct {
	mock *MockMetricsProvider
}

// NewMockMetricsProvider creates a new mock instance.
func NewMockMetricsProvider(ctrl *gomock.Controller) *MockMetricsProvider {
	mock := &MockMetricsProvider{ctrl: ctrl}
	mock.recorder = &MockMetricsProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsProvider) EXPECT() *MockMetricsProviderMockRecorder {
	return m.recorder
}

// CreateMetricURL mocks base method.
func (m *MockMetricsProvider) CreateMetricURL(mType constants.MetricType, mainURL, name string, value float64) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMetricURL", mType, mainURL, name, value)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateMetricURL indicates an expected call of CreateMetricURL.
func (mr *MockMetricsProviderMockRecorder) CreateMetricURL(mType, mainURL, name, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMetricURL", reflect.TypeOf((*MockMetricsProvider)(nil).CreateMetricURL), mType, mainURL, name, value)
}

// GetStorage mocks base method.
func (m *MockMetricsProvider) GetStorage() storage.MetricsRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorage")
	ret0, _ := ret[0].(storage.MetricsRepository)
	return ret0
}

// GetStorage indicates an expected call of GetStorage.
func (mr *MockMetricsProviderMockRecorder) GetStorage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorage", reflect.TypeOf((*MockMetricsProvider)(nil).GetStorage))
}

// SendMetricsHTTP mocks base method.
func (m *MockMetricsProvider) SendMetricsHTTP(reportInterval int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendMetricsHTTP", reportInterval)
}

// SendMetricsHTTP indicates an expected call of SendMetricsHTTP.
func (mr *MockMetricsProviderMockRecorder) SendMetricsHTTP(reportInterval interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMetricsHTTP", reflect.TypeOf((*MockMetricsProvider)(nil).SendMetricsHTTP), reportInterval)
}
