//go:build tools

package main

import (
	_ "github.com/favadi/protoc-go-inject-tag"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/wolfogre/gtag/cmd/gtag"
	_ "go.uber.org/mock/mockgen"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
