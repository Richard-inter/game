# Game Backend Services

A microservices-based game backend system built with Go, featuring claw machine, gacha machine, whack-a-mole games and player management services with real-time communication capabilities.

## 🏗️ Architecture

This project implements a distributed microservices architecture with the following components:

- **API Service** - REST API gateway (Port 8080)
- **WebSocket Service** - Real-time WebSocket connections (Port 8081)
- **TCP Service** - TCP socket connections (Port 8082)
- **RPC Services** - gRPC microservices:
  - ClawMachine Service (Port 9091)
  - ClawMachine Runtime Service (Port 9092)
  - GachaMachine Service (Port 9093)
  - GachaMachine Runtime Service (Port 9094)
  - WhackAMole Service (Port 9095)
  - WhackAMole Runtime Service (Port 9096)
  - Player Service (Port 9097)
- **Infrastructure** - MySQL database, etcd service discovery, Redis caching

## 🚀 Features

- **Microservices Architecture** - Scalable, independent services
- **Multi-Game Support** - Claw machine, gacha machine, and whack-a-mole games
- **Real-time Communication** - WebSocket and TCP support
- **Service Discovery** - etcd-based service registration and discovery
- **Database Integration** - MySQL with GORM ORM
- **Caching Layer** - Redis for performance optimization
- **Protocol Buffers** - Efficient inter-service communication
- **Docker Support** - Containerized deployment
- **Configuration Management** - YAML-based configuration
- **Leaderboard System** - Automated ranking and recalculation

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
make run-clawmachine-runtime # ClawMachine Runtime service (Port 9092)
make run-gachamachine  # GachaMachine RPC service (Port 9093)
make run-gachamachine-runtime # GachaMachine Runtime service (Port 9094)
make run-whackamole    # WhackAMole RPC service (Port 9095)
make run-whackamole-runtime # WhackAMole Runtime service (Port 9096)
make run-player        # Player RPC service (Port 9097)
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
make build-clawmachine-runtime
make build-gachamachine
make build-gachamachine-runtime
make build-whackamole
make build-whackamole-runtime
make build-player
make build-recalculate-leaderboard
```

## 📁 Project Structure

```
game/
├── cmd/                    # Service entry points
│   ├── api-service/        # REST API gateway
│   ├── websocket-service/  # WebSocket handler
│   ├── tcp-service/        # TCP socket handler
│   ├── recalculate-leaderboard/ # Leaderboard utility
│   └── rpc/               # gRPC microservices
│       ├── rpc-clawmachine-service/
│       ├── rpc-clawmachine-runtime-service/
│       ├── rpc-gachamachine-service/
│       ├── rpc-gachamachine-runtime-service/
│       ├── rpc-whackamole-service/
│       ├── rpc-whackamole-runtime-service/
│       └── rpc-player-service/
├── internal/              # Private application code
│   ├── cache/            # Redis caching layer
│   ├── config/           # Configuration management
│   ├── db/               # Database connections and seed data
│   ├── discovery/        # Service discovery (etcd)
│   ├── domain/           # Business logic models
│   ├── registry/         # Service registry
│   ├── repository/       # Data access layer
│   ├── service/          # Business services
│   ├── transport/        # Transport layer (HTTP, gRPC, WebSocket, TCP)
│   └── worker/           # Background workers
├── pkg/                  # Public library code
│   ├── common/           # Common utilities
│   ├── logger/           # Logging utilities
│   └── protocol/         # Protocol definitions (Proto, WebSocket)
│       ├── clawMachine/
│       ├── clawMachine_Websocket/
│       ├── gachaMachine/
│       ├── gachaMachine_Websocket/
│       ├── player/
│       ├── whackAMole/
│       └── whackAMole_Websocket/
├── config/               # Configuration files
├── scripts/              # Development scripts
├── docs/                 # API documentation
├── clawMachine.proto     # Protocol buffer definitions
├── simulate_gachamachine_client.go  # Gacha machine simulation client
├── test_whackamole_runtime.go      # Whack-a-mole runtime test
└── websocket_request.go    # WebSocket request utilities
```

## ⚙️ Configuration

Services use YAML configuration files in the `config/` directory:

- `api-service.yaml` - API service configuration
- `websocket-service.yaml` - WebSocket service configuration
- `tcp-service.yaml` - TCP service configuration
- `rpc-clawmachine-service.yaml` - ClawMachine RPC service configuration
- `rpc-clawmachine-runtime-service.yaml` - ClawMachine Runtime service configuration
- `rpc-gachamachine-service.yaml` - GachaMachine RPC service configuration
- `rpc-gachamachine-runtime-service.yaml` - GachaMachine Runtime service configuration
- `rpc-whackamole-service.yaml` - WhackAMole RPC service configuration
- `rpc-whackamole-runtime-service.yaml` - WhackAMole Runtime service configuration
- `rpc-player-service.yaml` - Player RPC service configuration
- `shared.yaml` - Shared configuration across services
- `config.yaml` - Global configuration

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

Test WhackAMole runtime service:
```bash
make test-whackamole-runtime
```

Simulate GachaMachine client:
```bash
go run simulate_gachamachine_client.go
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

### ClawMachine Runtime Service (Port 9092)
Runtime service for claw machine game sessions and real-time interactions.

### GachaMachine Service (Port 9093)
gRPC service managing gacha machine game logic and prize distribution.

### GachaMachine Runtime Service (Port 9094)
Runtime service for gacha machine game sessions and real-time interactions.

### WhackAMole Service (Port 9095)
gRPC service managing whack-a-mole game logic and scoring.

### WhackAMole Runtime Service (Port 9096)
Runtime service for whack-a-mole game sessions and real-time interactions.

### Player Service (Port 9097)
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

- **Gacha Machine Management**
  - `GET /gachamachines` - List gacha machines
  - `POST /gachamachines` - Create gacha machine
  - `GET /gachamachines/{id}` - Get gacha machine details
  - `PUT /gachamachines/{id}` - Update gacha machine
  - `DELETE /gachamachines/{id}` - Delete gacha machine

- **WhackAMole Management**
  - `GET /whackamoles` - List whack-a-mole games
  - `POST /whackamoles` - Create whack-a-mole game
  - `GET /whackamoles/{id}` - Get whack-a-mole details
  - `PUT /whackamoles/{id}` - Update whack-a-mole game
  - `DELETE /whackamoles/{id}` - Delete whack-a-mole game

- **Leaderboard Management**
  - `GET /leaderboards` - List leaderboards
  - `GET /leaderboards/{gameType}` - Get game-specific leaderboard
  - `POST /leaderboards/recalculate` - Trigger leaderboard recalculation

## 🗄️ Database

The project uses MySQL 8.0 as the primary database with Redis for caching. The database schema includes:

- `players` - Player information and statistics
- `claw_machines` - Claw machine configuration and state
- `gacha_machines` - Gacha machine configuration and prizes
- `whackamoles` - Whack-a-mole game configuration
- `games` - Game sessions and results
- `leaderboards` - Player rankings and scores
- `game_sessions` - Active game session tracking

Database connection is managed through GORM ORM with connection pooling. Redis is used for caching game states, player sessions, and leaderboard data for improved performance.

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

1. **Port conflicts** - Ensure ports 8080-8082 and 9090-9097 are available
2. **Database connection** - Verify MySQL is running and credentials are correct
3. **Redis connection** - Check Redis server is accessible for caching
4. **Service discovery** - Check etcd is accessible if enabled
5. **Docker issues** - Ensure Docker daemon is running and ports are exposed
6. **Game runtime issues** - Verify runtime services are started before game services

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