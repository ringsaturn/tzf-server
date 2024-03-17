GOBASEPATH=$(shell go env var GOPATH | xargs)

.PHONY:build
build:
	go build

.PHONY:test
test:
	go test -json -race ./... -v -coverprofile=coverage.txt -covermode=atomic | tparse -all

.PHONY:bench
bench:
	go test -bench=. -benchmem ./...

.PHONY:cover
cover:
	go tool cover -html=coverage.txt -o=coverage.html

mock:
	mockgen -source=$(GOBASEPATH)/pkg/mod/github.com/tidwall/redcon@v1.6.2/redcon.go -destination="internal/redisserver/mock_redcon_test.go" -package=redisserver_test

install:
	go install go.uber.org/mock/mockgen@latest
	go install github.com/wolfogre/gtag/cmd/gtag@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/ringsaturn/protoc-gen-go-hertz@latest

.PHONY:pb
pb:
	buf build
	buf generate
	protoc-go-inject-tag -input="tzf/v1/*.pb.go" -remove_tag_comment
	go fmt ./...

fmt:
	find proto/v1 -iname *.proto | xargs clang-format -i --style=Google
	go fmt ./...
	go fix ./...

gtag:
	cd internal/config;gtag -types Config -tags flag .

gen:
	make pb
	make mock
	make gtag
