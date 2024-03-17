package main

import (
	"context"
	"fmt"
	"log"
)

func main() {
	ctx := context.Background()
	a, err := newApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(a.Start(ctx))
}
