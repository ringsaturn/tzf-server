package config

type Config struct {
	FinderType int
	DataPath   string

	HTTPAddr          string
	DisablePrintRoute bool

	HertzPrometheusHostPorts    string
	HertzPrometheusPath         string
	PrometheusEnableGoCollector bool
	PrometheusPath              string
	PrometheusHostPorts         string
}
