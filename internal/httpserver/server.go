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

func bindAPI(h *server.Hertz, srv v1.TZFServiceHTTPServer) {
	v1.RegisterTZFServiceHTTPServer(h, srv)
	v1.BindDefaultSwagger(h)
}

func NewServer(cfg *config.Config, srv v1.TZFServiceHTTPServer) *server.Hertz {
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

	// h.SetHTMLTemplate(template.Must(template.New("").ParseFS(f, "template/*.html")))

	// h.GET("/", Index)
	// h.GET("/api/v1/tz/geojson", GetTimezoneShape)

	bindAPI(h, srv)

	// webPageGroup := h.Group("/web")
	// webPageGroup.GET("/tz", GetTimezoneInfoPage)
	// webPageGroup.GET("/tzs", GetTimezonesInfoPage)
	// webPageGroup.GET("/tzs/all", GetAllSupportTimezoneNamesPage)
	// webPageGroup.GET("/tz/geojson/viewer", GetGeoJSONViewerForTimezone)
	// webPageGroup.GET("/click", GetClickPage)
	return h
}

var ProviderSet = wire.NewSet(NewServer)
