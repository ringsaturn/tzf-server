package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/app/server"
	hc "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/google/wire"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/ringsaturn/tzf-server/internal/config"
	v1 "github.com/ringsaturn/tzf-server/tzf/v1"
)

func NewServer(cfg *config.Config, srv v1.TZFServiceHTTPServer, w *WebHandler) *server.Hertz {
	opts := []hc.Option{
		server.WithDisablePrintRoute(cfg.DisablePrintRoute),
		server.WithHostPorts(cfg.HTTPAddr),
		server.WithTracer(
			prometheus.NewServerTracer(
				cfg.HertzPrometheusHostPorts,
				cfg.HertzPrometheusPath,
				prometheus.WithEnableGoCollector(cfg.PrometheusEnableGoCollector),
			),
		),
	}

	h := server.New(opts...)

	h.Use(
		recovery.Recovery(recovery.WithRecoveryHandler(
			func(ctx context.Context, c *app.RequestContext, err interface{}, stack []byte) {
				c.JSON(http.StatusInternalServerError, utils.H{
					"error": fmt.Sprintf("[Recovery] err=%v\nstack=%s", err, stack),
				})
			},
		)),
	)

	bindAPI(h, srv)
	BindWebPage(h, w)

	return h
}

var ProviderSet = wire.NewSet(NewServer, NewWebHandler)
