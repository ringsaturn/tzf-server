package app

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/wire"
	"github.com/ringsaturn/tzf-server/internal/redisserver"
)

type App struct {
	Hertz       *server.Hertz
	RedisServer *redisserver.Server
}

func NewApp(h *server.Hertz, rs *redisserver.Server) *App {
	return &App{
		Hertz:       h,
		RedisServer: rs,
	}
}

var ProviderSet = wire.NewSet(NewApp)
