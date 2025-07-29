.PHONY: run run-port run-all run-custom run-cli build build-cli build-all test fmt clean docker-build docker-run help

APP_NAME=bitchest
PKG=github.com/daniacca/$(APP_NAME)
PORT=7463
HOST=0.0.0.0

# Default target
.DEFAULT_GOAL := run

# Compile and run the server
run:
	go run cmd/server/main.go

# Run server on custom port
run-port:
	go run cmd/server/main.go -port $(PORT)

# Run server on all interfaces
run-host:
	go run cmd/server/main.go -host $(HOST)

# Run server with custom host and port
run-host-port:
	go run cmd/server/main.go -host $(HOST) -port $(PORT)

# Compile and run the CLI client
run-cli:
	go run cmd/cli/main.go

# Compile the executable
build:
	go build -o ./out/$(APP_NAME) cmd/server/main.go

# Compile the CLI client
build-cli:
	go build -o ./out/$(APP_NAME)-cli cmd/cli/main.go

# Build both server and CLI
build-all: build build-cli

# Run tests
test:
	go test ./...

# Check formatting and lint
fmt:
	go fmt ./...

# Clean binaries or temporary builds
clean:
	rm -f $(APP_NAME)
	rm -f ./out/$(APP_NAME)
	rm -f ./out/$(APP_NAME)-cli

# Create Docker image
docker-build:
	docker build -t $(APP_NAME):latest .

# Run Docker container
docker-run:
	docker run --rm -p 7463:7463 $(APP_NAME):latest

# Show server help
help: build
	./out/$(APP_NAME) -h
