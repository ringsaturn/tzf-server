package httpserver

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	v1 "github.com/ringsaturn/tzf-server/tzf/v1"
)

// BindAPI binds the API to the server.
// This function is called by the generated code.
func BindAPI(h *server.Hertz, srv v1.TZFServiceHTTPServer) {
	v1.RegisterTZFServiceHTTPServer(h, srv)
	v1.BindDefaultSwagger(h)
}
