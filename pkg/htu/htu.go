package htu

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
)

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	TRACE   = "TRACE"
	OPTIONS = "OPTIONS"
)

type ValidationFunc func(t *testing.T, resp *ut.ResponseRecorder)

func Simple(
	t *testing.T,
	s *server.Hertz,
	method string, url string, body io.Reader,
	validate ValidationFunc,
) {
	w := ut.PerformRequest(
		s.Engine, method, url,
		nil,
	)
	validate(t, w)
}

type ValidationFuncForJSONAPI func(t *testing.T, expectedResponse interface{})

func JSONAPI(
	t *testing.T,
	s *server.Hertz,
	method string, url string,
	expectedResponse interface{},
	validate ValidationFuncForJSONAPI,
) {
	w := ut.PerformRequest(
		s.Engine, method, url,
		nil,
	)

	if err := json.Unmarshal(w.Body.Bytes(), expectedResponse); err != nil {
		t.Error(err.Error())
	}

	validate(t, expectedResponse)
}
