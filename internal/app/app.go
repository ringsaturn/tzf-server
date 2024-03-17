package app

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ringsaturn/tzf-server/internal/config"
	"github.com/ringsaturn/tzf-server/internal/redisserver"
	"golang.org/x/sync/errgroup"
)

type App struct {
	cfg         *config.Config
	Hertz       *server.Hertz
	RedisServer *redisserver.Server
}

func NewApp(cfg *config.Config, h *server.Hertz, rs *redisserver.Server) *App {
	return &App{
		cfg:         cfg,
		Hertz:       h,
		RedisServer: rs,
	}
}

func (a *App) StartHTTPServer() error {
	return a.Hertz.Run()
}

func (a *App) StartRedisServer() error {
	return a.RedisServer.StartRedisServer(a.cfg.RedisAddr)
}

func (a *App) StartPrometheus(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.Handle(a.cfg.PrometheusPath, promhttp.Handler())

	app := &http.Server{
		Addr:           a.cfg.PrometheusHostPorts,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	go func() {
		<-ctx.Done()
		_ = app.Shutdown(context.Background())
	}()

	return app.ListenAndServe()
}

func (a *App) Start(ctx context.Context) error {
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(a.StartHTTPServer)
	g.Go(a.StartRedisServer)
	g.Go(func() error {
		return a.StartPrometheus(gCtx)
	})

	return g.Wait()
}

var ProviderSet = wire.NewSet(NewApp)
