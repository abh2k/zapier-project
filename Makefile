BINARY := bin/server

.PHONY: build test run

build:
	go build -o $(BINARY) ./cmd/server

test:
	go test ./...

run:
	go mod tidy
	go run ./cmd/server
