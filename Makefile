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
	mockgen -source=$(GOBASEPATH)/pkg/mod/github.com/tidwall/redcon@v1.6.2/redcon.go -destination="handler/mock_redcon_test.go" -package=handler_test

install:
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/ringsaturn/protoc-gen-go-hertz@latest

.PHONY:pb
pb:
	buf build
	buf generate
	protoc-go-inject-tag -input="proto/v1/*.pb.go" -remove_tag_comment

fmt:
	find proto/v1 -iname *.proto | xargs clang-format -i --style=Google
	go fmt ./...
	go fix ./...
