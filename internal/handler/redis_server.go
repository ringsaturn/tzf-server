package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

func redisGetTZCmd(conn redcon.Conn, cmd redcon.Command) {
	lng, lat, err := parseCoordinates(cmd)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	timezone_name := finder.GetTimezoneName(lng, lat)
	if timezone_name == "" {
		conn.WriteError("no tz found")
		return
	}
	conn.WriteString(timezone_name)
}

func redisGetTZsCmd(conn redcon.Conn, cmd redcon.Command) {
	lng, lat, err := parseCoordinates(cmd)
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	timezone_names, err := finder.GetTimezoneNames(lng, lat)
	if err != nil {
		conn.WriteError("no tz found")
		return
	}
	conn.WriteArray(len(timezone_names))
	for _, name := range timezone_names {
		conn.WriteBulkString(name)
	}
}

func redisDefaultCmd(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
}

func redisPingCmd(conn redcon.Conn, _ redcon.Command) { conn.WriteString("PONG") }

func redisQuitCmd(conn redcon.Conn, _ redcon.Command) { conn.WriteString("OK"); conn.Close() }

func RedisHandler(conn redcon.Conn, cmd redcon.Command) {
	inputCmd := strings.ToLower(string(cmd.Args[0]))
	timer := prometheus.NewTimer(redisServerCmdHistogram.WithLabelValues(inputCmd))
	defer timer.ObserveDuration()
	switch inputCmd {
	case "ping":
		redisPingCmd(conn, cmd)
	case "quit":
		redisQuitCmd(conn, cmd)
	case "get_tz":
		redisGetTZCmd(conn, cmd)
	case "get_tzs":
		redisGetTZsCmd(conn, cmd)
	default:
		redisDefaultCmd(conn, cmd)
	}
}

func StartRedisServer(addr string) error {
	err := redcon.ListenAndServe(addr,
		RedisHandler,
		func(conn redcon.Conn) bool { return true },
		func(conn redcon.Conn, err error) {},
	)
	return err
}
