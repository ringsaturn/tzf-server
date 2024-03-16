package v1

import (
	"context"
	_ "embed"

	app "github.com/cloudwego/hertz/pkg/app"
	server "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/swagger"
	swaggerFiles "github.com/swaggo/files"
)

//go:embed openapi.yaml
var openapiYAML []byte

func BindSwagger(h *server.Hertz, openAPIYAMLPath string, swaggerPath string) {
	h.GET(swaggerPath, swagger.WrapHandler(
		swaggerFiles.Handler,
		swagger.URL(openAPIYAMLPath),
	))

	h.GET(openAPIYAMLPath, func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Content-Type", "application/x-yaml")
		_, _ = ctx.Write(openapiYAML)
	})
}

func BindDefaultSwagger(h *server.Hertz) {
	BindSwagger(h, "/openapi.yaml", "/swagger/*any")
}
