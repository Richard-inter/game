# Individual RPC Services Guide

## 🎯 **Individual Service Architecture**

Instead of one combined game service, each RPC service now has its own `main.go` and can be started independently.

## 📁 **Service Structure**

```
cmd/
├── game-service/main.go          # Combined service (legacy)
├── api-service/main.go           # HTTP API service
├── websocket-service/main.go       # WebSocket service
├── tcp-service/main.go            # TCP service
├── clawmachine-service/main.go    # NEW: ClawMachine gRPC service
└── player-service/main.go         # NEW: Player gRPC service

config/
├── game-service.yaml             # Combined service config
├── api-service.yaml             # HTTP API config
├── websocket-service.yaml         # WebSocket config
├── tcp-service.yaml             # TCP config
├── clawmachine-service.yaml      # NEW: ClawMachine config
└── player-service.yaml          # NEW: Player config
```

## 🚀 **Service Ports**

| Service | Type | Port | Config File | Command |
|---------|------|------|-------------|----------|
| ClawMachine | gRPC | 9091 | `clawmachine-service.yaml` | `make run-clawmachine` |
| Player | gRPC | 9092 | `player-service.yaml` | `make run-player` |
| API | HTTP | 8080 | `api-service.yaml` | `make run-api` |
| WebSocket | WS | 8081 | `websocket-service.yaml` | `make run-websocket` |
| TCP | TCP | 8082 | `tcp-service.yaml` | `make run-tcp` |

## 🔧 **Build Commands**

### **Build All Services:**
```bash
make build
```

### **Build Individual Services:**
```bash
make build-clawmachine    # Build ClawMachine service
make build-player        # Build Player service
make build-api           # Build API service
make build-websocket     # Build WebSocket service
make build-tcp           # Build TCP service
```

## 🏃 **Run Commands**

### **Run Individual Services:**
```bash
make run-clawmachine    # Start ClawMachine service (port 9091)
make run-player        # Start Player service (port 9092)
make run-api           # Start API service (port 8080)
make run-websocket     # Start WebSocket service (port 8081)
make run-tcp           # Start TCP service (port 8082)
```

### **With Custom Config:**
```bash
CONFIG_PATH=config/custom-clawmachine.yaml make run-clawmachine
CONFIG_PATH=config/custom-player.yaml make run-player
```

## 📋 **Service Configuration**

### **ClawMachine Service (`config/clawmachine-service.yaml`):**
```yaml
service:
  name: "clawmachine-service"
  host: "0.0.0.0"
  port: 9091
  mode: "release"

grpc:
  host: "0.0.0.0"
  port: 9091
  reflection: true

shared:
  database: "shared.yaml"
  redis: "shared.yaml"
  logging: "shared.yaml"
  jwt: "shared.yaml"
  tracing: "shared.yaml"
```

### **Player Service (`config/player-service.yaml`):**
```yaml
service:
  name: "player-service"
  host: "0.0.0.0"
  port: 9092
  mode: "release"

grpc:
  host: "0.0.0.0"
  port: 9092
  reflection: true

shared:
  database: "shared.yaml"
  redis: "shared.yaml"
  logging: "shared.yaml"
  jwt: "shared.yaml"
  tracing: "shared.yaml"
```

## 🔄 **Development Workflow**

### **1. Start Infrastructure:**
```bash
# Start MySQL only
make docker-up-infra
```

### **2. Start Services Individually:**
```bash
# Terminal 1: Start ClawMachine service
make run-clawmachine

# Terminal 2: Start Player service  
make run-player

# Terminal 3: Start API service
make run-api

# Terminal 4: Start WebSocket service
make run-websocket

# Terminal 5: Start TCP service
make run-tcp
```

### **3. Development Benefits:**
- **Independent debugging** - Each service has its own logs
- **Selective testing** - Run only services you need
- **Resource management** - Control memory/CPU per service
- **Hot reload** - Restart individual services quickly
- **Isolation** - Service crashes don't affect others

## 🔍 **Service Testing**

### **Test Individual Services:**
```bash
# Test ClawMachine service
grpcurl -plaintext localhost:9091 list

# Test Player service
grpcurl -plaintext localhost:9092 list

# Test API service
curl http://localhost:8080/health

# Test WebSocket service
websocat ws://localhost:8081/ws

# Test TCP service
telnet localhost 8082
```

### **Service Health Checks:**
```bash
# All gRPC services support reflection
grpcurl -plaintext localhost:9091 describe
grpcurl -plaintext localhost:9092 describe

# Check service status
curl http://localhost:8080/health
```

## 🎯 **VS Code Debugging**

### **Launch Configurations:**
Each service can be debugged independently in VS Code:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch ClawMachine Service",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/clawmachine-service/main.go",
            "env": {
                "CONFIG_PATH": "config/clawmachine-service.yaml"
            }
        },
        {
            "name": "Launch Player Service",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/cmd/player-service/main.go",
            "env": {
                "CONFIG_PATH": "config/player-service.yaml"
            }
        }
    ]
}
```

## ✅ **Benefits**

### **Microservice Architecture:**
- **Independent deployment** - Deploy services separately
- **Scalable** - Scale individual services
- **Resilient** - Service isolation prevents cascading failures
- **Flexible** - Mix and match service combinations
- **Debuggable** - Isolated debugging per service

### **Development Workflow:**
- **Focused development** - Work on one service at a time
- **Parallel development** - Multiple developers on different services
- **Testing** - Unit test individual services
- **CI/CD** - Separate build/deploy pipelines
- **Monitoring** - Individual service metrics

You now have complete individual RPC services that can be started and managed independently! 🚀
