package server

import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	mock_storage "DevOpsMetricsProject/internal/storage/mocks"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type statusCheckStruct struct {
	name         string
	endpoint     string
	method       string
	wantedStatus int
}

func StatusCheckTest(t *testing.T, tests []*statusCheckStruct, dompserv *dompserver) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, dompserv)
			require.NotNil(t, dompserv.coreMux)
			require.NotNil(t, dompserv.coreStg)

			request := httptest.NewRequest(tt.method, tt.endpoint, nil)
			resp := httptest.NewRecorder()
			dompserv.coreMux.ServeHTTP(resp, request)
			result := resp.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.wantedStatus, result.StatusCode)
		})
	}
}

func Test_updateMetricHandler(t *testing.T) {

	tests := []*statusCheckStruct{
		{
			name:         "UpdateMetricHandler getting correct URL",
			endpoint:     "/update/gauge/testGaugeName/64.51",
			method:       http.MethodPost,
			wantedStatus: http.StatusOK,
		},
		{
			name:         "UpdateMetricHandler getting wrong HTTP Method",
			endpoint:     "/update/gauge/testGaugeName/64.51",
			method:       http.MethodGet,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "UpdateMetricHandler getting URL with empty metric name",
			endpoint:     "/update/gauge//64.51",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "UpdateMetricHandler getting URL with no metric name #1",
			endpoint:     "/update/gauge/",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "UpdateMetricHandler getting URL with no metric name #2",
			endpoint:     "/update/gauge",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
	}
	cfg := &configs.ServerConfig{SaveMode: constants.InMemoryMode}

	serv, err := CreateNewServer(cfg)

	assert.Nil(t, err)

	StatusCheckTest(t, tests, serv)
}

func Test_getMainPageHandler(t *testing.T) {

	tests := []*statusCheckStruct{
		{
			name:         "GetMainPageHandler main case",
			endpoint:     "/",
			method:       http.MethodGet,
			wantedStatus: http.StatusOK,
		},
	}
	serv, err := CreateNewServer(&configs.ServerConfig{SaveMode: constants.FileMode})
	assert.Nil(t, err)
	StatusCheckTest(t, tests, serv)
}

func Test_dompserver_GetMetricHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	testingStorage := mock_storage.NewMockMetricsRepository(mockCtrl)

	tests := []*statusCheckStruct{
		{
			name:         "GetMetricHandler main case Gauge type",
			endpoint:     "/value/gauge/testName",
			method:       http.MethodGet,
			wantedStatus: http.StatusOK,
		},
		{
			name:         "GetMetricHandler main case Counter type",
			endpoint:     "/value/counter/testName",
			method:       http.MethodGet,
			wantedStatus: http.StatusOK,
		},
	}
	serv, err := CreateNewServer(&configs.ServerConfig{SaveMode: constants.FileMode})
	assert.Nil(t, err)

	serv.coreStg = testingStorage

	testingStorage.EXPECT().GetMetricByName(constants.GaugeType, gomock.Any()).AnyTimes().Return(55.5, nil)
	testingStorage.EXPECT().GetMetricByName(constants.CounterType, gomock.Any()).AnyTimes().Return(float64(55), nil)

	StatusCheckTest(t, tests, serv)
}
