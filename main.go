package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ringsaturn/tzf-server/internal/handler"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

var (
	finderType                  = flag.Int("type", 0, "which finder to use Polygon(0) or Fuzzy(1)")
	dataPath                    = flag.String("path", "", "custom data")
	httpAddr                    = flag.String("http-addr", "0.0.0.0:8080", "HTTP Host&Port")
	redisAddr                   = flag.String("redis-addr", "localhost:6380", "Redis Server Host&Port")
	prometheusHostPorts         = flag.String("prometheus-host-port", "0.0.0.0:2112", "Prometheus Host&Port")
	prometheusPath              = flag.String("prometheus-path", "/metrics", "Prometheus Path")
	hertzPrometheusHostPorts    = flag.String("hertz-prometheus-host-port", "0.0.0.0:8090", "Hertz Prometheus Host&Port")
	hertzPrometheusPath         = flag.String("hertz-prometheus-path", "/hertz", "Hertz Prometheus Path")
	prometheusEnableGoCollector = flag.Bool("prometheus-enable-go-coll", true, "Enable Go Collector")
	disablePrintRoute           = flag.Bool("disable-print-route", false, "Disable Print Route")
	debug                       = flag.Bool("debug", false, "Enable debug mode")
)

func main() {
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
				*hertzPrometheusHostPorts,
				*hertzPrometheusPath,
				prometheus.WithEnableGoCollector(*prometheusEnableGoCollector),
			),
		),
	)

	rootCtx := context.Background()

	g, gCtx := errgroup.WithContext(rootCtx)

	g.Go(h.Run)

	g.Go(func() error {
		mux := http.NewServeMux()
		mux.Handle(*prometheusPath, promhttp.Handler())

		app := &http.Server{
			Addr:           *prometheusHostPorts,
			ReadTimeout:    30 * time.Second,
			WriteTimeout:   30 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        mux,
		}

		go func() {
			<-gCtx.Done()
			_ = app.Shutdown(context.Background())
		}()

		return app.ListenAndServe()
	})

	g.Go(func() error { return handler.StartRedisServer(*redisAddr) })

	err = g.Wait()
	if err != nil {
		logger.Error("error", zap.Error(err))
		os.Exit(1)
	}
}
