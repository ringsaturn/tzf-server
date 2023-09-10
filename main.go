package main

import (
	"context"
	"flag"

	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/ringsaturn/tzf-server/internal/handler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

func main() {
	finderType := flag.Int("type", 0, "which finder to use Polygon(0) or Fuzzy(1)")
	dataPath := flag.String("path", "", "custom data")
	httpAddr := flag.String("http-addr", "localhost:8080", "HTTP Host&Port")
	redisAddr := flag.String("redis-addr", "localhost:6380", "Redis Server Host&Port")
	prometheusHostPorts := flag.String("prometheus-host-port", "localhost:8090", "Prometheus Host&Port")
	prometheusPath := flag.String("prometheus-path", "/hertz", "Prometheus Path")
	prometheusEnableGoCollector := flag.Bool("prometheus-enable-go-coll", true, "Enable Go Collector")
	disablePrintRoute := flag.Bool("disable-print-route", false, "Disable Print Route")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	var logger *zap.Logger
	var err error
	if *debug {
		cfg := zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		logger, err = cfg.Build()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	logger.Info("starting")
	logger.Sugar().Infow("Get config",
		"debug", *debug,
		"type", *finderType,
		"path", *dataPath,
		"http-addr", *httpAddr,
		"redis-addr", *redisAddr,
		"prometheus-host-port", *prometheusHostPorts,
		"prometheus-path", prometheusPath,
		"prometheus-enable-go-coll", *prometheusEnableGoCollector,
		"disable-print-route", *disablePrintRoute,
	)

	switch *finderType {
	case 0:
		logger.Debug("Will use Polygon data")
	case 1:
		logger.Debug("Will use Fuzzy data")
	default:
		logger.Error("Unknown method, quit.")
		return
	}

	if *dataPath == "" {
		logger.Debug("Will use built-in tzf-rel data")
	} else {
		logger.Debug("Will use custom data")
	}

	h := handler.Setup(
		logger,
		&handler.SetupFinderOptions{
			FinderType:     handler.FinderType((*finderType)),
			CustomDataPath: *dataPath,
		},
		server.WithDisablePrintRoute(*disablePrintRoute),
		server.WithHostPorts(*httpAddr),
		server.WithTracer(
			prometheus.NewServerTracer(
				*prometheusHostPorts,
				*prometheusPath,
				prometheus.WithEnableGoCollector(*prometheusEnableGoCollector),
			),
		),
	)

	rootCtx := context.Background()

	g, _ := errgroup.WithContext(rootCtx)

	g.Go(h.Run)

	g.Go(func() error { return handler.StartRedisServer(*redisAddr) })

	panic(g.Wait())
}
