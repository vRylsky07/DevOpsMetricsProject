package server

import (
	"DevOpsMetricsProject/internal/constants"
	"net/http"
	"net/http/httptest"
	"testing"

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
		{
			name:         "UpdateMetricHandler getting URL with only rout #1",
			endpoint:     "/update/",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "UpdateMetricHandler getting URL with only rout #2",
			endpoint:     "/update",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
	}
	StatusCheckTest(t, tests, CreateNewServer())
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
	StatusCheckTest(t, tests, CreateNewServer())
}

func Test_dompserver_GetMetricHandler(t *testing.T) {
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

	testServ := CreateNewServer()
	testServ.coreStg.UpdateMetricByName(constants.RenewOperation, constants.GaugeType, "testName", 55.5)
	testServ.coreStg.UpdateMetricByName(constants.RenewOperation, constants.CounterType, "testName", 55)

	StatusCheckTest(t, tests, testServ)
}
