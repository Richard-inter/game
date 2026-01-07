# Development Guide

## Overview

This project supports both Docker-based deployment and local development using VS Code's launch configurations. For development, using VS Code launch configurations is recommended over Docker.

## Development Setup

### 1. Prerequisites

- Go 1.21+
- VS Code with Go extension
- Docker & Docker Compose (for dependencies only)
- MySQL 8.0+ (or use Docker)
- Redis 7+ (or use Docker)

### 2. Initial Setup

```bash
# Clone and setup
git clone <repository-url>
cd game

# Install dependencies
make deps

# Generate protocol files
make proto

# Start dependencies (MySQL, Redis)
docker-compose up -d mysql redis
```

### 3. VS Code Setup

#### Copy Launch Configuration

```bash
# Copy launch configuration to VS Code directory
mkdir -p .vscode
cp launch.json .vscode/launch.json

# Copy tasks configuration
cp tasks.json .vscode/tasks.json
```

#### Available Launch Configurations

**Individual Services:**
- `Launch Game Service (gRPC)` - Start game service on port 9090
- `Launch API Service (HTTP)` - Start API service on port 8080
- `Launch WebSocket Service` - Start WebSocket service on port 8081
- `Launch TCP Service` - Start TCP service on port 8082

**Debug Configurations:**
- `Debug Game Service` - Debug game service with breakpoints
- `Debug API Service` - Debug API service with breakpoints

**Compound Configurations:**
- `Launch All Services` - Start all 4 services simultaneously
- `Launch Core Services` - Start Game + API services only
- `Launch Real-time Services` - Start WebSocket + TCP services only

## Using VS Code Launch

### Method 1: Run and Debug Panel

1. Open VS Code
2. Go to **Run and Debug** panel (Ctrl+Shift+D)
3. Select configuration from dropdown
4. Press **F5** or click **Play** button

### Method 2: Command Palette

1. Press **Ctrl+Shift+P**
2. Type `Debug: Select and Start Debugging`
3. Choose desired service

### Method 3: Quick Launch

- **Ctrl+F5**: Quick launch last used configuration
- **Shift+F5**: Launch without debugging

## Development Workflow

### 1. Start Dependencies

```bash
# Start only required dependencies
docker-compose up -d mysql redis

# Or start all services in Docker
docker-compose up -d
```

### 2. Launch Services Individually

For focused development on a single service:

```bash
# Using VS Code Launch
# Select "Launch Game Service (gRPC)" and press F5

# Or using command line
make run-game
make run-api
make run-websocket
make run-tcp
```

### 3. Debug Services

For debugging with breakpoints:

```bash
# Using VS Code Debug
# Select "Debug Game Service" and press F5
# Set breakpoints in your code
# Use debug console and variables
```

### 4. Launch Multiple Services

For testing service interactions:

```bash
# Using VS Code Compound Launch
# Select "Launch All Services" and press F5

# Or using multiple terminals
make run-game &    # Terminal 1
make run-api &     # Terminal 2
make run-websocket & # Terminal 3
make run-tcp &      # Terminal 4
```

## VS Code Tasks

Available tasks (Ctrl+Shift+P → "Tasks: Run Task"):

- `Build Game Service` - Build only game service
- `Build API Service` - Build only API service
- `Build WebSocket Service` - Build only WebSocket service
- `Build TCP Service` - Build only TCP service
- `Build All Services` - Build all services
- `Start Dependencies (Docker)` - Start MySQL/Redis only
- `Stop Dependencies` - Stop all Docker services
- `Generate Protocol Buffers` - Regenerate .proto files
- `Run Tests` - Run all tests

## Environment Variables

Launch configurations include development environment variables. For different environments:

### Development Environment (Default)
```bash
GAME_GAME_SERVICE_PORT=9090
GAME_API_SERVICE_PORT=8080
GAME_WEBSOCKET_SERVICE_PORT=8081
GAME_TCP_SERVICE_PORT=8082
GAME_DATABASE_HOST=localhost
GAME_REDIS_HOST=localhost
```

### Production Environment
Create `.env` file or set environment variables:
```bash
GAME_GAME_SERVICE_PORT=9090
GAME_DATABASE_HOST=prod-db-server
GAME_REDIS_HOST=prod-redis-server
GAME_DATABASE_PASSWORD=secure-password
```

## Hot Reload

For development with hot reload:

### Using Air (Recommended)

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Create .air.toml for each service
# Or use the global air configuration

# Run with hot reload
air -c .air.toml ./cmd/game-service
```

### Using VS Code Go Extension

1. Install Go extension for VS Code
2. Enable `go.delveConfig` settings
3. Use debug configurations for hot reload

## Testing

### Unit Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific service tests
go test ./internal/service/...
go test ./cmd/game-service/...
```

### Integration Tests

```bash
# Start dependencies
docker-compose up -d mysql redis

# Run integration tests
go test -tags=integration ./tests/...

# Stop dependencies
docker-compose down
```

## Port Management

Development ports:
- **Game Service (gRPC)**: 9090
- **API Service (HTTP)**: 8080
- **WebSocket Service**: 8081
- **TCP Service**: 8082
- **MySQL**: 3306
- **Redis**: 6379

## Troubleshooting

### Port Conflicts

If ports are in use:

```bash
# Find what's using the port
lsof -i :9090
lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change ports in launch.json
```

### Database Connection Issues

```bash
# Check if MySQL is running
docker-compose ps mysql

# Check MySQL logs
docker-compose logs mysql

# Reset database
docker-compose down -v
docker-compose up -d mysql
```

### Build Issues

```bash
# Clean build artifacts
make clean

# Regenerate protocol files
make proto

# Rebuild
make build
```

## Best Practices

### Development

1. **Use individual services** for focused development
2. **Use debug configurations** for troubleshooting
3. **Start only needed dependencies** to save resources
4. **Use environment variables** for different environments
5. **Run tests frequently** during development

### Code Organization

1. **Keep services independent** - avoid direct imports between services
2. **Use shared configs** for common settings
3. **Follow the existing patterns** for new services
4. **Update launch.json** when adding new services
5. **Document environment variables** in service configs

### Git Workflow

1. **Ignore build artifacts** (already in .gitignore)
2. **Don't commit .vscode/settings** (personal preferences)
3. **Commit launch.json** for team consistency
4. **Use feature branches** for new development
5. **Test services independently** before integration

## Production Deployment

For production, use Docker Compose:

```bash
# Build and deploy all services
docker-compose up -d

# Or deploy individual services
docker-compose up -d game-service
docker-compose up -d api-service
```

The launch configurations are primarily for development, not production deployment.
