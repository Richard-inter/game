# RPC Services Reorganization Guide

## 🎯 **Reorganization Complete**

Successfully reorganized RPC services to have "rpc" naming and grouped them into a dedicated `rpc` folder.

## 📁 **New Structure**

### **Command Directory:**
```
cmd/
├── api-service/                    # HTTP API service
├── game-service/                   # Combined game service (legacy)
├── server/                        # Legacy server
├── tcp-service/                   # TCP service
├── websocket-service/              # WebSocket service
└── rpc/                           # NEW: RPC services folder
    ├── rpc-clawmachine-service/    # ClawMachine gRPC service
    │   └── main.go
    └── rpc-player-service/          # Player gRPC service
        └── main.go
```

### **Config Directory:**
```
config/
├── api-service.yaml              # HTTP API config
├── game-service.yaml            # Combined game service config
├── tcp-service.yaml             # TCP service config
├── websocket-service.yaml         # WebSocket config
├── rpc-clawmachine-service.yaml  # NEW: ClawMachine RPC config
├── rpc-player-service.yaml       # NEW: Player RPC config
└── shared.yaml                 # Shared configuration
```

## 🚀 **Updated Services**

### **1. RPC ClawMachine Service**
- **Path**: `cmd/rpc/rpc-clawmachine-service/main.go`
- **Config**: `config/rpc-clawmachine-service.yaml`
- **Port**: 9091
- **Service Name**: `rpc-clawmachine-service`

### **2. RPC Player Service**
- **Path**: `cmd/rpc/rpc-player-service/main.go`
- **Config**: `config/rpc-player-service.yaml`
- **Port**: 9092
- **Service Name**: `rpc-player-service`

## 🔧 **Updated Makefile Commands**

### **Build Commands:**
```bash
make build-clawmachine    # Build RPC ClawMachine service
make build-player        # Build RPC Player service
make build               # Build all services (including RPC services)
```

### **Run Commands:**
```bash
make run-clawmachine    # Run RPC ClawMachine service (port 9091)
make run-player        # Run RPC Player service (port 9092)
```

### **Generated Binaries:**
```bash
bin/
├── rpc-clawmachine-service    # ClawMachine RPC binary
├── rpc-player-service        # Player RPC binary
├── api-service             # HTTP API binary
├── websocket-service       # WebSocket binary
└── tcp-service            # TCP binary
```

## 📋 **Service Configuration**

### **RPC ClawMachine Service (`config/rpc-clawmachine-service.yaml`):**
```yaml
# RPC ClawMachine Service Configuration (gRPC)

service:
  name: "rpc-clawmachine-service"
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

### **RPC Player Service (`config/rpc-player-service.yaml`):**
```yaml
# RPC Player Service Configuration (gRPC)

service:
  name: "rpc-player-service"
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
make docker-up-infra
```

### **2. Start RPC Services:**
```bash
# Terminal 1: Start RPC ClawMachine service
make run-clawmachine

# Terminal 2: Start RPC Player service
make run-player
```

### **3. With Custom Config:**
```bash
CONFIG_PATH=config/custom-rpc-clawmachine.yaml make run-clawmachine
CONFIG_PATH=config/custom-rpc-player.yaml make run-player
```

## 🗂️ **File Changes Summary**

### **Renamed Files:**
- `cmd/clawmachine-service/` → `cmd/rpc/rpc-clawmachine-service/`
- `cmd/player-service/` → `cmd/rpc/rpc-player-service/`
- `config/clawmachine-service.yaml` → `config/rpc-clawmachine-service.yaml`
- `config/player-service.yaml` → `config/rpc-player-service.yaml`

### **Updated Files:**
- `cmd/rpc/rpc-clawmachine-service/main.go` - Updated config path
- `cmd/rpc/rpc-player-service/main.go` - Updated config path
- `Makefile` - Updated build/run commands

## ✅ **Benefits of Reorganization**

### **Clear Naming:**
- **RPC prefix** - Clearly identifies RPC services
- **Consistent naming** - All RPC services follow same pattern
- **Service grouping** - All RPC services in dedicated folder

### **Better Organization:**
- **Logical grouping** - RPC services together
- **Easier navigation** - Clear folder structure
- **Scalable** - Easy to add new RPC services

### **Development Workflow:**
- **Independent services** - Each RPC service standalone
- **Separate configs** - Individual service configuration
- **Focused development** - Work on specific RPC services
- **Isolated debugging** - Separate logs per service

## 🎯 **Service Overview**

| Service Type | Service Name | Port | Command | Config File |
|-------------|---------------|------|----------|-------------|
| RPC gRPC | rpc-clawmachine-service | 9091 | `make run-clawmachine` | `rpc-clawmachine-service.yaml` |
| RPC gRPC | rpc-player-service | 9092 | `make run-player` | `rpc-player-service.yaml` |
| HTTP API | api-service | 8080 | `make run-api` | `api-service.yaml` |
| WebSocket | websocket-service | 8081 | `make run-websocket` | `websocket-service.yaml` |
| TCP | tcp-service | 8082 | `make run-tcp` | `tcp-service.yaml` |

Your RPC services are now properly reorganized with clear naming and structure! 🚀
