package server

import (
	//"DevOpsMetricsProject/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mainPageHandler(t *testing.T) {

	tests := []struct {
		name         string
		endpoint     string
		method       string
		wantedStatus int
	}{
		{
			name:         "Test handler getting correct URL",
			endpoint:     "/update/gauge/testGaugeName/64.51",
			method:       http.MethodPost,
			wantedStatus: http.StatusOK,
		},
		{
			name:         "Test handler getting wrong HTTP Method",
			endpoint:     "/update/gauge/testGaugeName/64.51",
			method:       http.MethodGet,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "Test handler getting URL with empty metric name",
			endpoint:     "/update/gauge//64.51",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "Test handler getting URL with too much parameters",
			endpoint:     "/update/gauge/testGaugeName/64.51/test/test",
			method:       http.MethodPost,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "Test handler getting URL with no metric value",
			endpoint:     "/update/gauge/testGaugeName",
			method:       http.MethodPost,
			wantedStatus: http.StatusBadRequest,
		},
		{
			name:         "Test handler getting URL with no metric name #1",
			endpoint:     "/update/gauge/",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "Test handler getting URL with no metric name #2",
			endpoint:     "/update/gauge",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "Test handler getting URL with only rout #1",
			endpoint:     "/update/",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "Test handler getting URL with only rout #2",
			endpoint:     "/update",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
		{
			name:         "Test handler getting URL wrong rout",
			endpoint:     "/main",
			method:       http.MethodPost,
			wantedStatus: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.endpoint, nil)
			resp := httptest.NewRecorder()
			MainPageHandler(resp, request)
			result := resp.Result()
			defer result.Body.Close()
			assert.Equal(t, tt.wantedStatus, result.StatusCode)
		})
	}
}
