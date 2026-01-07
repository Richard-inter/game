# Configuration Guide

## Overview

The game microservices use individual configuration files for better clarity and independence. Each service has its own config file that imports shared configurations.

## Config Structure

```
config/
├── shared.yaml              # Shared configs (database, redis, logging, etc.)
├── game-service.yaml        # Game gRPC service config
├── api-service.yaml         # HTTP API service config
├── websocket-service.yaml    # WebSocket service config
└── tcp-service.yaml         # TCP service config
```

## Individual Service Configs

### Game Service (`config/game-service.yaml`)

```yaml
service:
  name: "game-service"
  host: "0.0.0.0"
  port: 9090
  mode: "release"

grpc:
  host: "0.0.0.0"
  port: 9090
  reflection: true

# Import shared configurations
shared:
  database: "shared.yaml"
  redis: "shared.yaml"
  logging: "shared.yaml"
  jwt: "shared.yaml"
  tracing: "shared.yaml"
```

### API Service (`config/api-service.yaml`)

```yaml
service:
  name: "api-service"
  host: "0.0.0.0"
  port: 8080
  mode: "debug"
  read_timeout: 30
  write_timeout: 30

cors:
  allowed_origins: ["*"]
  allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowed_headers: ["*"]

# Import shared configurations
shared:
  database: "shared.yaml"
  redis: "shared.yaml"
  logging: "shared.yaml"
  jwt: "shared.yaml"
  tracing: "shared.yaml"
```

### WebSocket Service (`config/websocket-service.yaml`)

```yaml
service:
  name: "websocket-service"
  host: "0.0.0.0"
  port: 8081
  path: "/ws"
  read_buffer_size: 1024
  write_buffer_size: 1024

websocket:
  host: "0.0.0.0"
  port: 8081
  path: "/ws"
  read_buffer_size: 1024
  write_buffer_size: 1024
  check_origin: true

# Import shared configurations
shared:
  redis: "shared.yaml"
  logging: "shared.yaml"
  tracing: "shared.yaml"
```

### TCP Service (`config/tcp-service.yaml`)

```yaml
service:
  name: "tcp-service"
  host: "0.0.0.0"
  port: 8082
  keep_alive: true
  read_timeout: 30
  write_timeout: 30

tcp:
  host: "0.0.0.0"
  port: 8082
  keep_alive: true
  read_timeout: 30
  write_timeout: 30

# Import shared configurations
shared:
  logging: "shared.yaml"
  tracing: "shared.yaml"
```

### Shared Configuration (`config/shared.yaml`)

```yaml
database:
  host: "localhost"
  port: 3306
  user: "root"
  password: ""
  name: "game"
  charset: "utf8mb4"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, text
  output: "stdout"  # stdout, file

jwt:
  secret: "your-secret-key-change-in-production"
  expiration_time: 86400  # 24 hours

tracing:
  enabled: false
  service_name: "game-microservices"
  jaeger_url: "http://localhost:14268/api/traces"
```

## Environment Variables

Each service uses its own environment variable prefix:

- **Game Service**: `GAME_GAME_SERVICE_*`
- **API Service**: `GAME_API_SERVICE_*`
- **WebSocket Service**: `GAME_WEBSOCKET_SERVICE_*`
- **TCP Service**: `GAME_TCP_SERVICE_*`

### Examples

```bash
# Game Service
export GAME_GAME_SERVICE_PORT=9090
export GAME_GAME_SERVICE_HOST=localhost

# API Service
export GAME_API_SERVICE_PORT=8080
export GAME_API_SERVICE_MODE=release

# Shared configs (work for all services)
export GAME_DATABASE_HOST=localhost
export GAME_REDIS_HOST=localhost
export GAME_LOGGING_LEVEL=debug
```

## Loading Configuration

Each service loads its configuration using:

```go
cfg, err := config.LoadServiceConfig("service-name")
```

This automatically:
1. Loads the service-specific config file
2. Merges shared configurations
3. Applies environment variable overrides
4. Provides helper methods for addresses

### Helper Methods

```go
// Get service-specific address
addr := cfg.GetServiceAddr()  // "0.0.0.0:9090"

// Get protocol-specific addresses
grpcAddr := cfg.GetGRPCAddr()
wsAddr := cfg.GetWebSocketAddr()
tcpAddr := cfg.GetTCPAddr()
```

## Benefits

✅ **Clear Separation**: Each service has its own configuration  
✅ **No Port Conflicts**: Services define their own ports  
✅ **Independent Deployment**: Services can be configured separately  
✅ **Shared Resources**: Common configs are shared and maintained once  
✅ **Environment Override**: Each service can have its own environment variables  
✅ **Easy Maintenance**: Changes to one service don't affect others  

## Adding New Services

When creating a new microservice:

1. Create `config/new-service.yaml`
2. Define service-specific settings
3. Import required shared configs
4. Use `config.LoadServiceConfig("new-service")` in main.go
5. Update Docker Compose with new service configuration

## Docker Configuration

Each service mounts its specific config:

```yaml
new-service:
  volumes:
    - ./config/new-service.yaml:/app/config/new-service.yaml:ro
    - ./config/shared.yaml:/app/config/shared.yaml:ro
  environment:
    - GAME_NEW_SERVICE_PORT=xxxx
```

## Best Practices

1. **Service-Specific Settings**: Keep service-specific settings in the service config
2. **Shared Resources**: Use shared.yaml for common resources (database, redis)
3. **Environment Override**: Use environment variables for deployment-specific values
4. **Port Management**: Ensure each service uses a unique port
5. **Security**: Never commit secrets to config files, use environment variables
