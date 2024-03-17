package config

import (
	"flag"

	"github.com/google/wire"
)

type Config struct {
	FinderType                  int    `flag:"type"`
	DataPath                    string `flag:"path"`
	HTTPAddr                    string `flag:"http-addr"`
	RedisAddr                   string `flag:"redis-addr"`
	PrometheusHostPorts         string `flag:"prometheus-host-port"`
	PrometheusPath              string `flag:"prometheus-path"`
	HertzPrometheusHostPorts    string `flag:"hertz-prometheus-host-port"`
	HertzPrometheusPath         string `flag:"hertz-prometheus-path"`
	PrometheusEnableGoCollector bool   `flag:"prometheus-enable-go-coll"`
	DisablePrintRoute           bool   `flag:"disable-print-route"`
}

func NewConfigFromArgs() *Config {
	cfg := &Config{}

	tags := cfg.TagsFlag()

	flag.IntVar(&cfg.FinderType, tags.FinderType, 0, "which finder to use Polygon(0) or Fuzzy(1)")
	flag.StringVar(&cfg.DataPath, tags.DataPath, "", "custom data")
	flag.StringVar(&cfg.HTTPAddr, tags.HTTPAddr, "0.0.0.0:8080", "HTTP Host&Port")
	flag.StringVar(&cfg.RedisAddr, tags.RedisAddr, "localhost:6380", "Redis Server Host&Port")
	flag.StringVar(&cfg.PrometheusHostPorts, tags.PrometheusHostPorts, "0.0.0.0:2112", "Prometheus Host&Port")
	flag.StringVar(&cfg.PrometheusPath, tags.PrometheusPath, "/metrics", "Prometheus Path")
	flag.StringVar(&cfg.HertzPrometheusHostPorts, tags.HertzPrometheusHostPorts, "0.0.0.0:8090", "Hertz Prometheus Host&Port")
	flag.StringVar(&cfg.HertzPrometheusPath, tags.HertzPrometheusPath, "/hertz", "Hertz Prometheus Path")
	flag.BoolVar(&cfg.PrometheusEnableGoCollector, tags.PrometheusEnableGoCollector, true, "Enable Go Collector")
	flag.BoolVar(&cfg.DisablePrintRoute, tags.DisablePrintRoute, false, "Disable Print Route")

	flag.Parse()

	return cfg
}

var ProviderSet = wire.NewSet(NewConfigFromArgs)
