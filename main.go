package main

import (
	"context"
	"flag"
	"strings"

	"github.com/ringsaturn/tzf-server/handler"
	"github.com/tidwall/redcon"
	"golang.org/x/sync/errgroup"
)

func main() {
	finderType := flag.Int("type", 0, "which finder to use Polygon(0) or Fuzzy(1)")
	dataPath := flag.String("path", "", "custom data")
	flag.Parse()

	h := handler.Setup(&handler.SetupFinderOptions{
		FinderType:     handler.FinderType((*finderType)),
		CustomDataPath: *dataPath,
	})

	g, _ := errgroup.WithContext(context.Background())

	g.Go(h.Run)

	g.Go(func() error {
		err := redcon.ListenAndServe(":6380",
			func(conn redcon.Conn, cmd redcon.Command) {
				switch strings.ToLower(string(cmd.Args[0])) {
				default:
					conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
				case "ping":
					conn.WriteString("PONG")
				case "quit":
					conn.WriteString("OK")
					conn.Close()
				case "get_tz":
					handler.RedisGetTZCmd(conn, cmd)
				case "get_tzs":
					handler.RedisGetTZsCmd(conn, cmd)
				}
			},
			func(conn redcon.Conn) bool { return true },
			func(conn redcon.Conn, err error) {},
		)
		return err
	})

	panic(g.Wait())
}
