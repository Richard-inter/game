# Game Microservices Boilerplate

A comprehensive Go microservices boilerplate for gaming applications with flexible transport layers supporting HTTP API, gRPC, WebSocket, and TCP protocols. Built with Protocol Buffers for gRPC communication and following clean architecture principles.

## Features

### 🚀 Multi-Service Architecture
- **Game Service** (gRPC): Core game logic and player management
- **API Service** (HTTP): RESTful endpoints for web clients
- **WebSocket Service** (WebSocket): Real-time bidirectional communication
- **TCP Service** (TCP Socket): Raw TCP connections for custom protocols

### 📦 Protocol Support
- **Protocol Buffers**: gRPC service definitions and serialization
- **FlatBuffers**: High-performance serialization (optional)
- **JSON**: Standard HTTP API format

### 🏗️ Architecture
- **Clean Architecture**: Separation of concerns with domain-driven design
- **Microservices Ready**: Each service can run independently
- **Configuration Management**: Viper-based configuration with environment support
- **Database Integration**: GORM with MySQL support
- **Caching**: Redis integration
- **Logging**: Structured logging with Logrus
- **Tracing**: OpenTelemetry support with Jaeger

## Project Structure

```
game/
├── cmd/
│   ├── game-service/        # gRPC game service (port 9090)
│   ├── api-service/         # HTTP REST API service (port 8080)
│   ├── websocket-service/   # WebSocket service (port 8081)
│   └── tcp-service/        # TCP socket service (port 8082)
├── internal/
│   ├── config/            # Configuration management
│   ├── domain/            # Domain models and interfaces
│   ├── repository/        # Data access layer
│   ├── service/          # Business logic layer
│   └── transport/        # Transport layer implementations
├── pkg/
│   ├── logger/           # Logging utilities
│   ├── protocol/        # Protocol definitions (ProtoBuf, FlatBuffers)
│   └── common/          # Common utilities
├── config/              # Configuration files
├── data/               # Data storage (MySQL init scripts)
├── deployments/        # Deployment configurations
├── docs/               # Documentation
└── scripts/            # Utility scripts
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 7+
- Protocol Buffers compiler (`protoc`)

### Installation

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd game
   make init
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Generate protocol files**:
   ```bash
   make proto
   make flatbuffers  # Optional
   ```

4. **Start development environment**:
   ```bash
   docker-compose up -d mysql redis
   ```

### Running Services

#### Individual Services
Run each service independently:

```bash
# Game Service (gRPC)
make run-game
# or
./bin/game-service

# API Service (HTTP)
make run-api
# or
./bin/api-service

# WebSocket Service
make run-websocket
# or
./bin/websocket-service

# TCP Service
make run-tcp
# or
./bin/tcp-service
```

#### All Services
Build all services:
```bash
make build
```

## Service Endpoints

### Game Service (gRPC) - Port 9090

**Available RPC Methods:**
- `CreateGame(CreateGameRequest) -> CreateGameResponse`
- `GetGame(GetGameRequest) -> GetGameResponse`
- `ListGames(ListGamesRequest) -> ListGamesResponse`
- `JoinGame(JoinGameRequest) -> JoinGameResponse`

- `CreatePlayer(CreatePlayerRequest) -> CreatePlayerResponse`
- `GetPlayer(GetPlayerRequest) -> GetPlayerResponse`
- `ListPlayers(ListPlayersRequest) -> ListPlayersResponse`

### API Service (HTTP) - Port 8080

#### Health Check
```bash
GET /health
```

#### Games
```bash
GET    /api/v1/games          # List games
GET    /api/v1/games/:id      # Get game
POST   /api/v1/games          # Create game
PUT    /api/v1/games/:id      # Update game
DELETE /api/v1/games/:id      # Delete game
```

#### Players
```bash
GET    /api/v1/players        # List players
GET    /api/v1/players/:id    # Get player
POST   /api/v1/players        # Create player
PUT    /api/v1/players/:id    # Update player
DELETE /api/v1/players/:id    # Delete player
```

### WebSocket Service - Port 8081

Connect to WebSocket:
```bash
ws://localhost:8081/ws
```

### TCP Service - Port 8082

Connect to TCP server:
```bash
telnet localhost 8082
```

## Configuration

Configuration is managed through `config/config.yaml` and environment variables. Environment variables use the `GAME_` prefix:

```bash
export GAME_SERVER_PORT=8080
export GAME_DATABASE_HOST=localhost
export GAME_REDIS_HOST=localhost
export GAME_GRPC_PORT=9090
export GAME_WEBSOCKET_PORT=8081
export GAME_TCP_PORT=8082
```

## Development

### Available Commands

```bash
# Build commands
make build           # Build all services
make build-game      # Build game service only
make build-api       # Build API service only
make build-websocket # Build WebSocket service only
make build-tcp       # Build TCP service only

# Run commands
make run-game        # Run game service
make run-api         # Run API service
make run-websocket   # Run WebSocket service
make run-tcp         # Run TCP service

# Other commands
make test            # Run tests
make test-coverage   # Run tests with coverage
make clean           # Clean build artifacts
make generate        # Generate code
make proto           # Generate protobuf files
make flatbuffers     # Generate FlatBuffer files
make fmt             # Format code
make lint            # Lint code
make dev             # Start development mode with hot reload
```

### Protocol Generation

#### Protocol Buffers (gRPC)
```bash
protoc --go_out=. --go-grpc_out=. pkg/protocol/*.proto
```

#### FlatBuffers
```bash
flatc --go -o pkg/protocol pkg/protocol/*.fbs
```

## Using Docker Compose

Start the entire stack:
```bash
docker-compose up -d
```

This will start:
- Game Service (gRPC:9090)
- API Service (HTTP:8080)
- WebSocket Service (WebSocket:8081)
- TCP Service (TCP:8082)
- MySQL Database (3306)
- Redis Cache (6379)
- Jaeger Tracing (16686)
- Prometheus Metrics (9091)
- Grafana Dashboard (3000)

## Service Communication

### Inter-Service Communication
Services can communicate via:
- **gRPC**: High-performance internal communication
- **HTTP REST**: For external integrations
- **Message Queues**: Redis pub/sub for async communication
- **Shared Database**: For data persistence

### Example gRPC Client
```go
import (
    "google.golang.org/grpc"
    pb "github.com/1nterdigital/game/pkg/protocol"
)

conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

client := pb.NewGameServiceClient(conn)
resp, err := client.CreateGame(context.Background(), &pb.CreateGameRequest{
    Name:        "My Game",
    Description: "A fun game",
    MaxPlayers:  10,
})
```

## Monitoring

### Health Checks
- Game Service: gRPC health check
- API Service: `GET /health`
- WebSocket Service: `GET /health`
- TCP Service: Connection test

### Metrics
- Prometheus: `http://localhost:9091/metrics`
- Grafana: `http://localhost:3000` (admin/admin)

### Tracing
- Jaeger UI: `http://localhost:16686`

## Architecture Decisions

### Microservices Design
Each service is designed to be:
- **Single Responsibility**: Focused on a specific domain
- **Independent**: Can be deployed and scaled separately
- **Resilient**: Failure in one service doesn't affect others
- **Technology Agnostic**: Can use different databases/technologies per service

### Protocol Choice
- **gRPC**: For internal service communication (high performance)
- **HTTP**: For external APIs and web clients
- **WebSocket**: For real-time features
- **TCP**: For custom protocols and legacy systems

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make fmt && make lint && make test`
6. Submit a pull request

## License

This project is licensed under the MIT License.
