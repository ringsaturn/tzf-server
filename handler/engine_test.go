package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ringsaturn/gtu"
	"github.com/ringsaturn/tzf-server/handler"
	"github.com/stretchr/testify/assert"
)

var (
	engine      = handler.Setup(nil)
	fuzzyEngine = handler.Setup(&handler.SetupFinderOptions{FinderType: handler.FuzzyFinder})
)

func TestEngine(t *testing.T) {
	type args struct {
		t        *testing.T
		engine   *gin.Engine
		method   string
		url      string
		body     io.Reader
		validate gtu.ValidationFunc
		options  []gtu.RequestOption
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Index",
			args: args{
				t:      t,
				engine: engine,
				method: gtu.GET,
				url:    "/",
				body:   nil,
				validate: func(t *testing.T, httpResp *httptest.ResponseRecorder) {
					assert.Equal(t, httpResp.Result().StatusCode, http.StatusTemporaryRedirect)
				},
			},
		},
		{
			name: "Ping",
			args: args{
				t:      t,
				engine: engine,
				method: gtu.GET,
				url:    "/ping",
				body:   nil,
				validate: func(t *testing.T, httpResp *httptest.ResponseRecorder) {
					assert.Equal(t, httpResp.Result().StatusCode, http.StatusOK)
				},
			},
		},
		{
			name: "GetTimezoneName",
			args: args{
				t:      t,
				engine: engine,
				method: gtu.GET,
				url:    "/api/v1/tz?lng=116.3883&lat=39.9289",
				body:   nil,
				validate: func(t *testing.T, httpResp *httptest.ResponseRecorder) {
					assert.Equal(t, httpResp.Result().StatusCode, http.StatusOK)
					resp := map[string]string{}
					err := json.Unmarshal(httpResp.Body.Bytes(), &resp)
					if err != nil {
						t.Fatalf(err.Error())
					}
					assert.Equal(t, resp["timezone"], "Asia/Shanghai")
				},
			},
		},
		{
			name: "FuzzyEngine-GetTimezoneName",
			args: args{
				t:      t,
				engine: fuzzyEngine,
				method: gtu.GET,
				url:    "/api/v1/tz?lng=116.3883&lat=39.9289",
				body:   nil,
				validate: func(t *testing.T, httpResp *httptest.ResponseRecorder) {
					assert.Equal(t, httpResp.Result().StatusCode, http.StatusOK)
					resp := map[string]string{}
					err := json.Unmarshal(httpResp.Body.Bytes(), &resp)
					if err != nil {
						t.Fatalf(err.Error())
					}
					assert.Equal(t, resp["timezone"], "Asia/Shanghai")
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gtu.Simple(tt.args.t, tt.args.engine, tt.args.method, tt.args.url, tt.args.body, tt.args.validate, tt.args.options...)
		})
	}
}
