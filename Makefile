.PHONY: build run test clean proto generate docker-build docker-run

# Variables
APP_NAME := game
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION := $(shell go version | awk '{print $$3}')

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GoVersion=$(GO_VERSION)"

# Default target
all: build

# Build the application
build:
	@echo "Building all services..."
	mkdir -p bin
	go build $(LDFLAGS) -o bin/game-service ./cmd/game-service
	go build $(LDFLAGS) -o bin/api-service ./cmd/api-service
	go build $(LDFLAGS) -o bin/websocket-service ./cmd/websocket-service
	go build $(LDFLAGS) -o bin/tcp-service ./cmd/tcp-service
	go build $(LDFLAGS) -o bin/rpc-clawmachine-service ./cmd/rpc/rpc-clawmachine-service
	go build $(LDFLAGS) -o bin/rpc-player-service ./cmd/rpc/rpc-player-service

# Build individual services
build-game:
	@echo "Building game service..."
	go build $(LDFLAGS) -o bin/game-service ./cmd/game-service

build-api:
	@echo "Building API service..."
	go build $(LDFLAGS) -o bin/api-service ./cmd/api-service

build-websocket:
	@echo "Building WebSocket service..."
	go build $(LDFLAGS) -o bin/websocket-service ./cmd/websocket-service

build-tcp:
	@echo "Building TCP service..."
	go build $(LDFLAGS) -o bin/tcp-service ./cmd/tcp-service

build-clawmachine:
	@echo "Building RPC ClawMachine service..."
	go build $(LDFLAGS) -o bin/rpc-clawmachine-service ./cmd/rpc/rpc-clawmachine-service

build-player:
	@echo "Building RPC Player service..."
	go build $(LDFLAGS) -o bin/rpc-player-service ./cmd/rpc/rpc-player-service

# Run individual services
run-game:
	@echo "Running game service..."
	go run $(LDFLAGS) ./cmd/game-service

run-api:
	@echo "Running API service..."
	go run $(LDFLAGS) ./cmd/api-service

run-websocket:
	@echo "Running WebSocket service..."
	go run $(LDFLAGS) ./cmd/websocket-service

run-tcp:
	@echo "Running TCP service..."
	go run $(LDFLAGS) ./cmd/tcp-service

run-clawmachine:
	@echo "Running RPC ClawMachine service..."
	go run $(LDFLAGS) ./cmd/rpc/rpc-clawmachine-service

run-player:
	@echo "Running RPC Player service..."
	go run $(LDFLAGS) ./cmd/rpc/rpc-player-service

# Run the application (all services)
run:
	@echo "Running all services..."
	@echo "Starting game service on port 9090..."
	@echo "Starting API service on port 8080..."
	@echo "Starting WebSocket service on port 8081..."
	@echo "Starting TCP service on port 8082..."
	@echo "Use 'make run-game', 'make run-api', etc. to run individual services"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate code (proto, flatbuffers, etc.)
generate:
	@echo "Generating code..."
	go generate ./...

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	find pkg/protocol -name "*.proto" | xargs protoc \
		-I pkg/protocol \
		--go_out=paths=source_relative:pkg/protocol \
		--go-grpc_out=paths=source_relative:pkg/protocol

# Generate flatbuffer files
flatbuffers:
	@echo "Generating flatbuffer files..."
	flatc --go -o pkg/protocol pkg/protocol/*.fbs

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker build -t game-services:latest .

docker-build-dev:
	@echo "Building development Docker image..."
	docker build -f Dockerfile.dev -t game-services:dev .

docker-up:
	@echo "Starting all services..."
	docker-compose up -d

docker-up-infra:
	@echo "Starting infrastructure only..."
	docker-compose -f docker-compose.mysql.yml up -d

docker-down:
	@echo "Stopping all services..."
	docker-compose down

docker-down-infra:
	@echo "Stopping infrastructure only..."
	docker-compose -f docker-compose.mysql.yml down

docker-logs:
	@echo "Showing logs..."
	docker-compose logs -f

docker-logs-infra:
	@echo "Showing infrastructure logs..."
	docker-compose -f docker-compose.infrastructure.yml logs -f

docker-logs-dev:
	@echo "Showing development service logs..."
	docker-compose -f docker-compose.dev.yml logs -f

# Development workflow commands
dev-setup:
	@echo "Setting up development environment..."
	make docker-up-infra
	@echo "Infrastructure started. Now use VS Code launch configurations to start services."

dev-full:
	@echo "Starting full development environment..."
	make docker-up-infra
	@echo "Infrastructure started. Use VS Code launch configurations for microservices."

dev-clean:
	@echo "Cleaning development environment..."
	make docker-down-infra
	docker system prune -f
	# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Lint code with fix
lint-fix:
	@echo "Linting and fixing code..."
	golangci-lint run --fix

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Development mode (watch and rebuild)
dev:
	@echo "Starting development mode..."
	air -c .air.toml

# Initialize project
init:
	@echo "Initializing project..."
	go mod init github.com/Richard-inter/game
	go mod tidy
	mkdir -p bin cmd internal pkg config data docs deployments scripts

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  generate       - Generate code"
	@echo "  proto          - Generate protobuf files"
	@echo "  flatbuffers    - Generate flatbuffer files"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  deps           - Install dependencies"
	@echo "  lint           - Lint code"
	@echo "  fmt            - Format code"
	@echo "  dev            - Start development mode"
	@echo "  init           - Initialize project"
	@echo "  help           - Show this help message"
