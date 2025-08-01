# Build targets
build: clean
	go build -o ./out/bitchest cmd/server/main.go

build-cli: clean
	go build -o ./out/bitchest-cli cmd/cli/main.go

build-all: build build-cli

# Test targets
test:
	go test ./...

test-verbose:
	go test -v ./...

# Coverage targets
coverage:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

coverage-summary:
	go test -cover ./...

coverage-clean:
	rm -f coverage.out coverage.html

# Run targets
run:
	go run cmd/server/main.go

run-host:
	go run cmd/server/main.go -host $(HOST)

run-port:
	go run cmd/server/main.go -port $(PORT)

run-host-port:
	go run cmd/server/main.go -host $(HOST) -port $(PORT)

run-cli:
	go run cmd/cli/main.go

# Docker targets
docker-build:
	docker build -t bitchest:latest .

docker-run:
	docker run --rm -p 7463:7463 bitchest:latest

# Help target
help:
	@echo "Available targets:"
	@echo "  build          - Build server binary"
	@echo "  build-cli      - Build CLI client binary"
	@echo "  build-all      - Build both server and CLI"
	@echo "  test           - Run tests"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  coverage       - Generate coverage report"
	@echo "  coverage-summary - Show coverage summary"
	@echo "  coverage-clean - Clean coverage files"
	@echo "  run            - Run server on localhost:7463"
	@echo "  run-host       - Run server on 0.0.0.0:7463"
	@echo "  run-port       - Run server on localhost:6379"
	@echo "  run-host-port  - Run server on 0.0.0.0:6379"
	@echo "  run-cli        - Run CLI client"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  clean          - Clean build artifacts"
	@echo "  help           - Show this help message"

# Clean target
clean:
	rm -f ./out/bitchest ./out/bitchest-cli

# Variables
APP_NAME = bitchest
PKG = github.com/daniacca/$(APP_NAME)
PORT = 6379
HOST = 0.0.0.0

# Phony targets
.PHONY: build build-cli build-all test test-verbose coverage coverage-summary coverage-clean run run-host run-port run-host-port run-cli docker-build docker-run help clean
