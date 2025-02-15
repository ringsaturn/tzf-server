GOBASEPATH=$(shell go env var GOPATH | xargs)

help:     ## Show this help.
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'


.PHONY:build
build:  ## Build tzf server
	go build

.PHONY:test
test:  ## Run test
	go test -json -race ./... -v -coverprofile=coverage.txt -covermode=atomic | tparse -all

.PHONY:bench
bench:  ## Run benchmark
	go test -bench=. -benchmem ./...

.PHONY:cover
cover:  ## Generate coverage report
	go tool cover -html=coverage.txt -o=coverage.html

mock:  ## Generate mock
	mockgen -source=$(GOBASEPATH)/pkg/mod/github.com/tidwall/redcon@v1.6.2/redcon.go -destination="internal/redisserver/mock_redcon_test.go" -package=redisserver_test

install:  ## Install tools
	go install github.com/mfridman/tparse@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/wolfogre/gtag/cmd/gtag@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/ringsaturn/protoc-gen-go-hertz@latest

bump-buf-dep:
	buf dep update

.PHONY:pb
pb:  ## Generate protobuf
	buf build
	buf generate
	protoc-go-inject-tag -input="gen/go/tzf_server/v1/*.pb.go" -remove_tag_comment
	cp internal/misc/api_hertz_swagger.go.txt gen/openapi/api_hertz_swagger.go
	go fmt ./...

fmt:  ## Format code
	find proto/v1 -iname *.proto | xargs clang-format -i --style=Google
	go fmt ./...
	go fix ./...

gtag:  ## Generate gtag
	cd internal/config;gtag -types Config -tags flag .

gen:  ## Generate code
	make pb
	make mock
	make gtag
