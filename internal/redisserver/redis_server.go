package redisserver

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/wire"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/ringsaturn/tzf"
	"github.com/tidwall/redcon"
)

var (
	redisServerCmdHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tzf_server_redis_server_cmd_histogram",
			Help:    "",
			Buckets: []float64{0.001, 0.002, 0.003, 0.004, 0.005, 0.006, 0.007, 0.008, 0.009, 0.01, 0.1},
		},
		[]string{"cmd"},
	)
)

func parseCoordinates(cmd redcon.Command) (float64, float64, error) {
	if len(cmd.Args) != 3 {
		return 0, 0, errors.New("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
	}
	lng, err := strconv.ParseFloat(string(cmd.Args[1]), 64)
	if err != nil {
		return 0, 0, err
	}
	lat, err := strconv.ParseFloat(string(cmd.Args[2]), 64)
	if err != nil {
		return 0, 0, err
	}
	return lng, lat, nil
}

type Server struct {
	f tzf.F
}

func NewServer(f tzf.F) *Server {
	return &Server{f: f}
}

var ProviderSet = wire.NewSet(NewServer)

func (s *Server) redisGetTZCmd(conn redcon.Conn, cmd redcon.Command) {
	lng, lat, err := parseCoordinates(cmd)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	timezone_name := s.f.GetTimezoneName(lng, lat)
	if timezone_name == "" {
		conn.WriteError("no tz found")
		return
	}
	conn.WriteString(timezone_name)
}

func (s *Server) redisGetTZsCmd(conn redcon.Conn, cmd redcon.Command) {
	lng, lat, err := parseCoordinates(cmd)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	timezone_names, err := s.f.GetTimezoneNames(lng, lat)
	if err != nil {
		conn.WriteError("no tz found")
		return
	}
	conn.WriteArray(len(timezone_names))
	for _, name := range timezone_names {
		conn.WriteBulkString(name)
	}
}

func (s *Server) Handler(conn redcon.Conn, cmd redcon.Command) {
	inputCmd := strings.ToLower(string(cmd.Args[0]))
	timer := prometheus.NewTimer(redisServerCmdHistogram.WithLabelValues(inputCmd))
	defer timer.ObserveDuration()
	switch inputCmd {
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "get_tz":
		s.redisGetTZCmd(conn, cmd)
	case "get_tzs":
		s.redisGetTZsCmd(conn, cmd)
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	}
}

func (s *Server) StartRedisServer(addr string) error {
	err := redcon.ListenAndServe(addr,
		s.Handler,
		func(conn redcon.Conn) bool { return true },
		func(conn redcon.Conn, err error) {},
	)
	return err
}
