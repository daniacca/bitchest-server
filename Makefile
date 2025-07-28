APP_NAME=bitchest
PKG=github.com/tuo-username/$(APP_NAME)

# Default target
.DEFAULT_GOAL := run

# Compile and run the server
run:
	go run cmd/server/main.go

# Compile the executable
build:
	go build -o ./out/$(APP_NAME) cmd/server/main.go

# Run tests
test:
	go test ./...

# Check formatting and lint
fmt:
	go fmt ./...

# Clean binaries or temporary builds
clean:
	rm -f $(APP_NAME)

# Create Docker image
docker-build:
	docker build -t $(APP_NAME):latest .

# Run Docker container
docker-run:
	docker run --rm -p 7463:7463 $(APP_NAME):latest
