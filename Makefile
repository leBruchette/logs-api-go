.PHONY: all start test

all: start

start:
	go run main/server.go

test:
	gotest -v ./...