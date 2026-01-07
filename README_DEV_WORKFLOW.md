# Development Workflow Guide

## Overview

This project uses a **hybrid development approach**:
- **Docker Compose** for infrastructure services (MySQL, Redis, etc.)
- **VS Code Launch Configurations** for microservices (API, Game, WebSocket, TCP)

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   VS Code       │    │  Docker Compose │    │  Your Local    │
│                 │    │                 │    │  Machine       │
│ ┌─────────────┐ │    │ ┌─────────────┐ │    │                 │
│ │ Launch      │ │    │ │ MySQL       │ │    │                 │
│ │ Configs    │ │    │ │ Redis       │ │    │                 │
│ └─────────────┘ │    │ │ Jaeger      │ │    │                 │
│                 │    │ │ Prometheus  │ │    │                 │
│ ┌─────────────┐ │    │ │ Grafana     │ │    │                 │
│ │ Game Svc    │◄──►│ └─────────────┘ │◄──►│                 │
│ │ API Svc     │ │    │                 │    │                 │
│ │ WebSocket   │ │    │                 │    │                 │
│ │ TCP Svc     │ │    │                 │    │                 │
│ └─────────────┘ │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### 1. Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd game

# Install Go dependencies
make deps

# Generate protocol files
make proto

# Setup VS Code configurations
mkdir -p .vscode
cp launch.json .vscode/launch.json
cp tasks.json .vscode/tasks.json
```

### 2. Start Infrastructure

```bash
# Start only infrastructure services (MySQL, Redis)
make docker-up-infra

# Or include observability (Jaeger, Prometheus, Grafana)
make docker-up-observability
```

### 3. Launch Microservices

**Method 1: VS Code (Recommended)**
- Open VS Code
- Go to **Run & Debug** panel (Ctrl+Shift+D)
- Select service and press **F5**

**Method 2: Command Line**
```bash
# Individual services
make run-game      # Game gRPC service
make run-api       # API HTTP service
make run-websocket # WebSocket service
make run-tcp       # TCP service
```

## 📋 Available Commands

### Infrastructure Commands
```bash
make docker-up-infra      # Start MySQL, Redis
make docker-up-observability  # Start + Jaeger, Prometheus, Grafana
make docker-down-infra     # Stop infrastructure
make docker-logs-infra    # View infrastructure logs
```

### Development Commands
```bash
make dev-setup           # Start infra + show instructions
make dev-full            # Start infra (observability included)
make dev-clean           # Stop everything + cleanup
```

### Service Commands
```bash
make build-game          # Build game service
make build-api           # Build API service
make build-websocket     # Build WebSocket service
make build-tcp           # Build TCP service
make build               # Build all services
```

## 🎯 Development Workflows

### Workflow 1: Focused Development

**When working on a single service:**

```bash
# 1. Start infrastructure
make docker-up-infra

# 2. Launch your service in VS Code
# Select "Debug Game Service" + F5

# 3. Code with breakpoints and hot reload
# Make changes → they reload automatically
```

### Workflow 2: Integration Testing

**When testing service interactions:**

```bash
# 1. Start infrastructure
make docker-up-infra

# 2. Launch multiple services
# Use "Launch All Services" in VS Code
# Or open multiple terminals:
make run-game &     # Terminal 1
make run-api &      # Terminal 2
make run-websocket & # Terminal 3
```

### Workflow 3: Full Stack Development

**When needing complete environment:**

```bash
# 1. Start everything
make dev-full

# 2. All services available:
# MySQL: localhost:3306
# Redis: localhost:6379
# Game gRPC: localhost:9090
# API HTTP: localhost:8080
# WebSocket: localhost:8081
# TCP: localhost:8082
# Jaeger: http://localhost:16686
# Prometheus: http://localhost:9091
# Grafana: http://localhost:3000
```

## 🔧 Configuration

### Infrastructure Configuration

**`docker-compose.infrastructure.yml`** contains:
- MySQL database
- Redis cache
- Jaeger tracing (optional)
- Prometheus metrics (optional)
- Grafana dashboards (optional)

### Service Configuration

**Each service uses its own config file:**
- `config/game-service.yaml`
- `config/api-service.yaml`
- `config/websocket-service.yaml`
- `config/tcp-service.yaml`
- `config/shared.yaml` (common settings)

### Environment Variables

**Infrastructure:**
```bash
# Database connection
GAME_DATABASE_HOST=localhost
GAME_DATABASE_PORT=3306
GAME_DATABASE_USER=root
GAME_DATABASE_PASSWORD=

# Redis connection
GAME_REDIS_HOST=localhost
GAME_REDIS_PORT=6379
```

**Services (via launch.json):**
```bash
# Service-specific ports
GAME_GAME_SERVICE_PORT=9090
GAME_API_SERVICE_PORT=8080
GAME_WEBSOCKET_SERVICE_PORT=8081
GAME_TCP_SERVICE_PORT=8082
```

## 🐛 Debugging

### VS Code Debugging

1. **Set breakpoints** in your Go code
2. **Select debug configuration** (e.g., "Debug Game Service")
3. **Press F5** to start with debugger
4. **Use debug console** for inspection
5. **Step through code** with F10/F11

### Service-Specific Debugging

**Game Service (gRPC):**
```bash
# Test with grpcurl
grpcurl -plaintext localhost:9090 list

# Test specific service
grpcurl -plaintext -d '{"name":"test"}' localhost:9090 game.GameService/CreateGame
```

**API Service (HTTP):**
```bash
# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/games
curl -X POST http://localhost:8080/api/v1/games -H "Content-Type: application/json" -d '{"name":"test"}'
```

**WebSocket Service:**
```bash
# Test with websocat
websocat ws://localhost:8081/ws

# Or use browser WebSocket client
# Connect to ws://localhost:8081/ws
```

**TCP Service:**
```bash
# Test with telnet
telnet localhost 8082

# Or with nc
nc localhost 8082
```

## 📊 Monitoring

### Infrastructure Health

```bash
# Check infrastructure status
docker-compose -f docker-compose.infrastructure.yml ps

# View logs
make docker-logs-infra

# Access services
# MySQL: mysql://localhost:3306
# Redis: redis://localhost:6379
# Jaeger: http://localhost:16686
# Prometheus: http://localhost:9091
# Grafana: http://localhost:3000 (admin/admin)
```

### Application Logs

**VS Code Debug Console:**
- Real-time logs during debugging
- Structured logging with logrus
- Error stack traces

**Command Line Logs:**
```bash
# Individual service logs
./bin/game-service
./bin/api-service
./bin/websocket-service
./bin/tcp-service
```

## 🔄 Hot Reload

### During Development

**VS Code Auto-Rebuild:**
- Go extension automatically rebuilds on save
- Debug session restarts automatically
- Fast feedback loop

**Manual Rebuild:**
```bash
# Rebuild and restart service
make build-game
./bin/game-service
```

## 🧪 Testing

### Unit Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Test specific package
go test ./internal/service/...
```

### Integration Tests

```bash
# Start infrastructure first
make docker-up-infra

# Run integration tests
go test -tags=integration ./tests/...

# Clean up
make docker-down-infra
```

## 🚀 Deployment

### Development Deployment

**For testing your changes:**

```bash
# Option 1: Local binaries
make build
./bin/game-service &
./bin/api-service &

# Option 2: Docker with hot reload
make docker-up-dev

# Option 3: Full Docker stack
make docker-up-all
```

### Production Deployment

**Use the production Docker Compose:**

```bash
# Deploy all services
docker-compose up -d

# Deploy specific services
docker-compose up -d game-service api-service
```

## 📝 Best Practices

### Development

1. **Use VS Code launch configs** for microservices
2. **Use Docker only for infrastructure**
3. **Start only needed services** to save resources
4. **Use debug configurations** for troubleshooting
5. **Test service interactions** with multiple launches

### Code Organization

1. **Keep services independent** - no direct imports
2. **Use shared configs** for common settings
3. **Follow existing patterns** for new services
4. **Update launch.json** when adding services
5. **Document environment variables**

### Git Workflow

1. **Ignore build artifacts** (in .gitignore)
2. **Commit launch.json** for team consistency
3. **Use feature branches** for new development
4. **Test locally** before pushing
5. **Document breaking changes**

## 🔍 Troubleshooting

### Port Conflicts

```bash
# Check what's using ports
lsof -i :9090  # Game service
lsof -i :8080  # API service
lsof -i :8081  # WebSocket
lsof -i :8082  # TCP

# Kill processes if needed
kill -9 <PID>
```

### Database Issues

```bash
# Check MySQL status
docker-compose -f docker-compose.infrastructure.yml ps mysql

# Reset database
make docker-down-infra
make docker-up-infra

# Access MySQL directly
docker-compose -f docker-compose.infrastructure.yml exec mysql mysql -u root -p
```

### Service Won't Start

```bash
# Check service logs in VS Code debug console
# Or run manually to see errors
./bin/game-service

# Check configuration
cat config/game-service.yaml
cat config/shared.yaml
```

This hybrid approach gives you the **best of both worlds**: fast, debuggable local development with reliable, isolated infrastructure services.
