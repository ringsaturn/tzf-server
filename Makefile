.PHONY:build
build:
	go build

.PHONY:test
test:
	go test -json -race ./... -v -coverprofile=coverage.out | tparse -all

.PHONY:cover
cover:
	go tool cover -html=coverage.out -o=coverage.html
