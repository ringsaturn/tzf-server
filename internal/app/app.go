package app

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/ringsaturn/tzf-server/internal/config"
	"github.com/ringsaturn/tzf-server/internal/redisserver"
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

var ProviderSet = wire.NewSet(NewApp)
