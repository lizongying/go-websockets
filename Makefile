.PHONY: all client server

all: client server

client:
	go vet ./cmd/client
	go build -ldflags "-s -w" -o  ./assets/client  ./cmd/client

server:
	go vet ./cmd/server
	go build -ldflags "-s -w" -o  ./assets/server  ./cmd/server