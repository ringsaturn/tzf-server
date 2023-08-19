package main

import (
	"context"
	"flag"

	"github.com/ringsaturn/tzf-server/internal/handler"
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

	g.Go(handler.StartRedisServer)

	panic(g.Wait())
}
