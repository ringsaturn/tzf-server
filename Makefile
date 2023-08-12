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
	go tool cover -html=coverage.out -o=coverage.html

gen:
	mockgen -source=$(GOBASEPATH)/pkg/mod/github.com/tidwall/redcon@v1.6.2/redcon.go -destination="handler/mock_redcon_test.go" -package=handler_test
