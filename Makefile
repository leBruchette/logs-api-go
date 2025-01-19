.PHONY: all start test

all: start

verify-deps:
	@command -v docker-compose >/dev/null 2>&1 || { echo >&2 "docker-compose is not installed. Aborting."; exit 1; }
	@docker info >/dev/null 2>&1 || { echo >&2 "Docker is not running. Please start Docker"; exit 1; }
	@echo "All dependencies are installed and Docker is running."

install-gotest:
	go get -u github.com/rakyll/gotest

test: install-gotest
	gotest -v ./...

start:
	go run main/server.go

docker-build:
	docker-compose build

docker-run:
	docker-compose up