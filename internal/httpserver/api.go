package httpserver

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	v1 "github.com/ringsaturn/tzf-server/gen/go/tzf_server/v1"
	apiV1Openapi "github.com/ringsaturn/tzf-server/gen/openapi"
)

// BindAPI binds the API to the server.
// This function is called by the generated code.
func BindAPI(h *server.Hertz, srv v1.TZFServiceHTTPServer) {
	v1.RegisterTZFServiceHTTPServer(h, srv)
	apiV1Openapi.BindDefaultSwagger(h)
}
