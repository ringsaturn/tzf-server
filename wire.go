//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/ringsaturn/tzf-server/internal/app"
	"github.com/ringsaturn/tzf-server/internal/config"
	"github.com/ringsaturn/tzf-server/internal/finder"
	"github.com/ringsaturn/tzf-server/internal/httpserver"
	"github.com/ringsaturn/tzf-server/internal/redisserver"
	"github.com/ringsaturn/tzf-server/internal/wraps"
)

func newApp(ctx context.Context) (*app.App, error) {
	panic(wire.Build(
		config.ProviderSet,
		redisserver.ProviderSet,
		finder.ProviderSet,
		wraps.HTTPProviderSet,
		app.ProviderSet,
		httpserver.ProviderSet,
	))
}
