# Game Microservices Boilerplate

A comprehensive Go microservices boilerplate for gaming applications with flexible transport layers supporting HTTP API, gRPC, WebSocket, and TCP protocols. Built with FlatBuffers for efficient serialization and following clean architecture principles.

## Features

### 🚀 Multi-Transport Architecture
- **HTTP API**: RESTful endpoints with Gin framework
- **gRPC**: High-performance RPC communication
- **WebSocket**: Real-time bidirectional communication
- **TCP Socket**: Raw TCP connections for custom protocols

### 📦 Protocol Support
- **FlatBuffers**: High-performance serialization
- **Protocol Buffers**: gRPC integration
- **JSON**: Standard HTTP API format

### 🏗️ Architecture
- **Clean Architecture**: Separation of concerns with domain-driven design
- **Microservices Ready**: Modular structure for distributed systems
- **Configuration Management**: Viper-based configuration with environment support
- **Database Integration**: GORM with MySQL support
- **Caching**: Redis integration
- **Logging**: Structured logging with Logrus
- **Tracing**: OpenTelemetry support with Jaeger

### 🛠️ Development Tools
- **Docker & Docker Compose**: Containerized deployment
- **Makefile**: Common development tasks
- **Hot Reload**: Air support for development
- **Testing**: Unit test framework

## Project Structure

```
game/
├── cmd/
│   └── server/           # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── domain/           # Domain models and interfaces
│   └── transport/        # Transport layer implementations
│       ├── grpc/         # gRPC server
│       ├── http/         # HTTP API server
│       ├── tcp/          # TCP socket server
│       └── websocket/    # WebSocket server
├── pkg/
│   ├── logger/           # Logging utilities
│   └── protocol/         # Protocol definitions (FlatBuffers)
├── config/               # Configuration files
├── data/                 # Data storage (MySQL init scripts)
├── deployments/          # Deployment configurations
├── docs/                 # Documentation
└── scripts/              # Utility scripts
```

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- MySQL 8.0+
- Redis 7+
- FlatBuffers compiler (`flatc`)

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
   make flatbuffers
   make proto
   ```

4. **Start development environment**:
   ```bash
   docker-compose up -d mysql redis
   make run
   ```

### Using Docker Compose

Start the entire stack:
```bash
docker-compose up -d
```

This will start:
- Game Server (HTTP:8080, gRPC:9090, WebSocket:8081, TCP:8082)
- MySQL Database (3306)
- Redis Cache (6379)
- Jaeger Tracing (16686)
- Prometheus Metrics (9091)
- Grafana Dashboard (3000)

## API Endpoints

### HTTP API

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

### WebSocket

Connect to WebSocket:
```bash
ws://localhost:8081/ws
```

### TCP Socket

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
```

## Development

### Available Commands

```bash
make build          # Build the application
make run            # Run the application
make test           # Run tests
make test-coverage  # Run tests with coverage
make clean          # Clean build artifacts
make generate       # Generate code
make flatbuffers    # Generate FlatBuffer files
make proto          # Generate protobuf files
make fmt            # Format code
make lint           # Lint code
make dev            # Start development mode with hot reload
```

### Adding New Services

1. Create domain models in `internal/domain/`
2. Implement repositories in `internal/repository/`
3. Implement services in `internal/service/`
4. Add handlers in `internal/transport/{http,grpc,websocket,tcp}/`
5. Update protocol definitions in `pkg/protocol/`

### Protocol Generation

#### FlatBuffers
```bash
flatc --go -o pkg/protocol pkg/protocol/*.fbs
```

#### Protocol Buffers
```bash
protoc --go_out=. --go-grpc_out=. pkg/protocol/*.proto
```

## Deployment

### Docker

Build and run:
```bash
docker build -t game:latest .
docker run -p 8080:8080 -p 9090:9090 -p 8081:8081 -p 8082:8082 game:latest
```

### Kubernetes

Deployments are available in `deployments/k8s/`:

```bash
kubectl apply -f deployments/k8s/
```

## Monitoring

### Health Checks
- HTTP: `GET /health`

### Metrics
- Prometheus: `http://localhost:9091/metrics`
- Grafana: `http://localhost:3000` (admin/admin)

### Tracing
- Jaeger UI: `http://localhost:16686`

## Architecture Decisions

### Multi-Transport Support
The application supports multiple transport layers to handle different use cases:
- **HTTP API**: For web clients and external integrations
- **gRPC**: For internal microservice communication
- **WebSocket**: For real-time gaming features
- **TCP**: For custom protocols and high-performance scenarios

### FlatBuffers Integration
FlatBuffers are used for efficient serialization in real-time communication:
- Zero-copy deserialization
- Forward/backward compatibility
- Cross-language support

### Clean Architecture
The project follows clean architecture principles:
- Domain logic is independent of frameworks
- Infrastructure concerns are isolated
- Testability is prioritized

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run `make fmt && make lint && make test`
6. Submit a pull request

## License

This project is licensed under the MIT License.
