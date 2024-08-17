package server

/*
import (
	"DevOpsMetricsProject/internal/configs"
	"DevOpsMetricsProject/internal/constants"
	storage_custom_mocks "DevOpsMetricsProject/internal/storage/custommock"

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
	}
	cfg := &configs.ServerConfig{SaveMode: constants.FileMode}

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

	g := map[string]float64{"testGauge": 55.5}
	c := map[string]int{"testCounter": 55}

	testingStorage := &storage_custom_mocks.StorageMockCustom{Gauge: g, Counter: c}
	serv.coreStg = testingStorage

	StatusCheckTest(t, tests, serv)
}
*/
