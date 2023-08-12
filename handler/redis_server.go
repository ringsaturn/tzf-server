package handler

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/redcon"
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

func redisPingCmd(conn redcon.Conn, cmd redcon.Command) { conn.WriteString("PONG") }

func redisQuitCmd(conn redcon.Conn, cmd redcon.Command) { conn.WriteString("OK"); conn.Close() }

func RedisHandler(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
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

func StartRedisServer() error {
	err := redcon.ListenAndServe(":6380",
		RedisHandler,
		func(conn redcon.Conn) bool { return true },
		func(conn redcon.Conn, err error) {},
	)
	return err
}
