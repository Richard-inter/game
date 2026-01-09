# Game Backend Services

A microservices-based game backend system built with Go, featuring claw machine and player management services with real-time communication capabilities.

## 🏗️ Architecture

This project implements a distributed microservices architecture with the following components:

- **API Service** - REST API gateway (Port 8080)
- **WebSocket Service** - Real-time WebSocket connections (Port 8081)
- **TCP Service** - TCP socket connections (Port 8082)
- **RPC Services** - gRPC microservices:
  - ClawMachine Service (Port 9091)
  - Player Service (Port 9094)
- **Infrastructure** - MySQL database, etcd service discovery

## 🚀 Features

- **Microservices Architecture** - Scalable, independent services
- **Real-time Communication** - WebSocket and TCP support
- **Service Discovery** - etcd-based service registration and discovery
- **Database Integration** - MySQL with GORM ORM
- **Protocol Buffers** - Efficient inter-service communication
- **Docker Support** - Containerized deployment
- **Configuration Management** - YAML-based configuration

## 📋 Prerequisites

- Go 1.24+
- Docker & Docker Compose
- MySQL 8.0+
- etcd (optional, for service discovery)
- Protocol Buffers compiler (protoc)
- FlatBuffers compiler (flatc)

## 🛠️ Installation

1. Clone the repository:
```bash
git clone https://github.com/Richard-inter/game.git
cd game
```

2. Install dependencies:
```bash
make deps
```

3. Start infrastructure:
```bash
make docker-up-infra
```

## 🏃‍♂️ Quick Start

### Development Mode

Start all services individually:
```bash
# Start each service in separate terminals
make run-game          # Game service (Port 9090)
make run-api           # API service (Port 8080)
make run-websocket     # WebSocket service (Port 8081)
make run-tcp           # TCP service (Port 8082)
make run-clawmachine   # ClawMachine RPC service (Port 9091)
make run-player        # Player RPC service (Port 9094)
```

### Docker Mode

Start all services with Docker Compose:
```bash
make docker-up
```

### Build

Build all services:
```bash
make build
```

Build individual services:
```bash
make build-api
make build-websocket
make build-tcp
make build-clawmachine
make build-player
```

## 📁 Project Structure

```
game/
├── cmd/                    # Service entry points
│   ├── api-service/        # REST API gateway
│   ├── websocket-service/  # WebSocket handler
│   ├── tcp-service/        # TCP socket handler
│   └── rpc/               # gRPC microservices
│       ├── rpc-clawmachine-service/
│       └── rpc-player-service/
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── db/               # Database connections
│   ├── discovery/        # Service discovery (etcd)
│   ├── domain/           # Business logic models
│   ├── repository/       # Data access layer
│   ├── service/          # Business services
│   └── transport/        # Transport layer (HTTP, gRPC, WebSocket, TCP)
├── pkg/                  # Public library code
│   ├── common/           # Common utilities
│   ├── logger/           # Logging utilities
│   └── protocol/         # Protocol definitions (Proto, FlatBuffers)
├── config/               # Configuration files
├── scripts/              # Development scripts
```

## ⚙️ Configuration

Services use YAML configuration files in the `config/` directory:

- `api-service.yaml` - API service configuration
- `websocket-service.yaml` - WebSocket service configuration
- `game-service.yaml` - Game service configuration
- `rpc-clawmachine-service.yaml` - ClawMachine RPC service configuration
- `rpc-player-service.yaml` - Player RPC service configuration

### Environment Variables

- `CONFIG_PATH` - Path to configuration file (optional, defaults to service-specific config)

## 🔧 Development

### Code Generation

Generate protocol buffer files:
```bash
make proto
```

Generate FlatBuffer files:
```bash
make flatbuffers
```

Generate all code:
```bash
make generate
```

### Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

### Linting

Lint code:
```bash
make lint
```

Lint and fix:
```bash
make lint-fix
```

### Formatting

Format code:
```bash
make fmt
```

## 🐳 Docker

### Development

Build development image:
```bash
make docker-build-dev
```

### Production

Build production image:
```bash
make docker-build
```

Start services:
```bash
make docker-up
```

Stop services:
```bash
make docker-down
```

View logs:
```bash
make docker-logs
```

## 📊 Services

### API Service (Port 8080)
REST API gateway that handles HTTP requests and forwards them to appropriate microservices.

### WebSocket Service (Port 8081)
Handles real-time WebSocket connections for live game updates.

### TCP Service (Port 8082)
Manages TCP socket connections for low-latency communication.

### ClawMachine Service (Port 9091)
gRPC service managing claw machine game logic and state.

### Player Service (Port 9094)
gRPC service handling player data, authentication, and profiles.

## 🔌 API Endpoints

The API service exposes REST endpoints for:

- **Player Management**
  - `GET /players` - List players
  - `POST /players` - Create player
  - `GET /players/{id}` - Get player details
  - `PUT /players/{id}` - Update player
  - `DELETE /players/{id}` - Delete player

- **Claw Machine Management**
  - `GET /clawmachines` - List claw machines
  - `POST /clawmachines` - Create claw machine
  - `GET /clawmachines/{id}` - Get claw machine details
  - `PUT /clawmachines/{id}` - Update claw machine
  - `DELETE /clawmachines/{id}` - Delete claw machine

## 🗄️ Database

The project uses MySQL 8.0 as the primary database. The database schema includes:

- `players` - Player information and statistics
- `claw_machines` - Claw machine configuration and state
- `games` - Game sessions and results

Database connection is managed through GORM ORM with connection pooling.

## 🔍 Service Discovery

Services can be discovered using etcd. When service discovery is enabled:

1. Services register themselves with etcd on startup
2. Client services discover service endpoints dynamically
3. Load balancing and health checking are handled automatically

To disable service discovery, set `discovery.enabled: false` in configuration.

## 🧪 Testing

The project includes unit tests for business logic and integration tests for API endpoints. Tests are located alongside the source code.

## 📝 Logging

Structured logging using zap with high-performance JSON output:
- DEBUG
- INFO  
- WARN
- ERROR
- FATAL

Logs include service name, version, request correlation IDs, and structured fields for better observability. The logger supports both structured logging with zap.SugaredLogger and typed logging with zap.Logger for optimal performance.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run linting and tests
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Port conflicts** - Ensure ports 8080-8082 and 9090-9094 are available
2. **Database connection** - Verify MySQL is running and credentials are correct
3. **Service discovery** - Check etcd is accessible if enabled
4. **Docker issues** - Ensure Docker daemon is running and ports are exposed

### Health Checks

Each service exposes a health check endpoint:
- `GET /health` - Service health status

### Logs

View service logs for debugging:
```bash
# Docker logs
make docker-logs

# Individual service logs
docker-compose logs -f [service-name]
```

## 📞 Support

For support and questions:
- Create an issue in the GitHub repository
- Check the documentation in the `docs/` directory
- Review configuration examples in `config/`