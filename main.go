package main

import (
	"flag"

	"github.com/ringsaturn/tzf-server/handler"
)

func main() {
	finderType := flag.Int("type", 0, "which finder to use Polygon(0) or Fuzzy(1)")
	dataPath := flag.String("path", "", "custom data")
	flag.Parse()

	h := handler.Setup(&handler.SetupFinderOptions{
		FinderType:     handler.FinderType((*finderType)),
		CustomDataPath: *dataPath,
	})
	panic(h.Run())
}
