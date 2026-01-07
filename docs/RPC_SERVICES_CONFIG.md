# RPC Services Configuration Guide

## 🎯 **Updated Configuration Structure**

### **New `game-service.yaml` Configuration:**

```yaml
# Game Service Configuration (gRPC)

service:
  name: "game-service"
  host: "0.0.0.0"
  port: 9090
  mode: "release"

grpc:
  host: "0.0.0.0"
  port: 9090
  reflection: true

# RPC Services Configuration
rpc_services:
  claw_machine:
    enabled: true
    service_package: "clawMachine"
    service_name: "ClawMachineService"
    import_path: "github.com/Richard-inter/game/pkg/protocol/clawMachine"
    implementation_path: "github.com/Richard-inter/game/internal/service/rpc/clawMachine"
  
  player:
    enabled: true
    service_package: "player"
    service_name: "PlayerService"
    import_path: "github.com/Richard-inter/game/pkg/protocol/player"
    implementation_path: "github.com/Richard-inter/game/internal/service/rpc/player"

# Import shared configurations
shared:
  database: "shared.yaml"
  redis: "shared.yaml"
  logging: "shared.yaml"
  jwt: "shared.yaml"
  tracing: "shared.yaml"
```

## 🏗️ **Configuration Structure**

### **RPC Service Config Fields:**

| Field | Type | Description |
|--------|------|-------------|
| `enabled` | bool | Whether the service is active |
| `service_package` | string | Protocol package name |
| `service_name` | string | gRPC service name |
| `import_path` | string | Protocol import path |
| `implementation_path` | string | Service implementation path |

### **Available Services:**

#### **1. ClawMachine Service:**
- **Protocol**: `pkg/protocol/clawMachine/`
  - `clawMachine.proto`
  - `clawMachine.pb.go`
  - `clawMachine_grpc.pb.go`
- **Implementation**: `internal/service/rpc/clawMachine/`
  - `clawMachine_service.go`
- **Service**: `ClawMachineService`

#### **2. Player Service:**
- **Protocol**: `pkg/protocol/player/`
  - `player.proto`
  - `player.pb.go`
  - `player_grpc.pb.go`
- **Implementation**: `internal/service/rpc/player/`
  - `player_service.go`
- **Service**: `PlayerService`

## 🔧 **Updated Config Structures**

### **New Config Types:**

```go
type RPCServiceConfig struct {
    Enabled           bool   `mapstructure:"enabled"`
    ServicePackage    string `mapstructure:"service_package"`
    ServiceName      string `mapstructure:"service_name"`
    ImportPath        string `mapstructure:"import_path"`
    ImplementationPath string `mapstructure:"implementation_path"`
}

type RPCServicesConfig struct {
    ClawMachine RPCServiceConfig `mapstructure:"claw_machine"`
    Player      RPCServiceConfig `mapstructure:"player"`
}

type Config struct {
    Server      ServerConfig        `mapstructure:"server"`
    Database    DatabaseConfig      `mapstructure:"database"`
    Redis       RedisConfig         `mapstructure:"redis"`
    GRPC        GRPCConfig          `mapstructure:"grpc"`
    RPCServices RPCServicesConfig   `mapstructure:"rpc_services"`  // NEW!
    WebSocket   WebSocketConfig     `mapstructure:"websocket"`
    TCP         TCPConfig           `mapstructure:"tcp"`
    JWT         JWTConfig           `mapstructure:"jwt"`
    Logging     LoggingConfig       `mapstructure:"logging"`
    Tracing     TracingConfig       `mapstructure:"tracing"`
}
```

## 🚀 **Usage in Application**

### **Loading Configuration:**
```go
// Load service-specific configuration
configFile := os.Getenv("CONFIG_PATH")
if configFile == "" {
    configFile = "config/game-service.yaml"
}

cfg, err := config.LoadServiceConfigFromPath(configFile)
if err != nil {
    log.WithError(err).Fatal("Failed to load configuration")
}
```

### **Accessing RPC Service Config:**
```go
// Check if claw machine service is enabled
if cfg.RPCServices.ClawMachine.Enabled {
    log.Info("ClawMachine service is enabled")
    log.WithField("path", cfg.RPCServices.ClawMachine.ImportPath).Info("Protocol import path")
}

// Check if player service is enabled
if cfg.RPCServices.Player.Enabled {
    log.Info("Player service is enabled")
    log.WithField("path", cfg.RPCServices.Player.ImportPath).Info("Protocol import path")
}
```

## 🔄 **Dynamic Service Loading**

### **RPC Service Loader:**
Created `internal/config/rpc_loader.go` for dynamic service loading:

```go
loader := config.NewRPCServiceLoader(logger)
err := loader.LoadServices(grpcServer, cfg.RPCServices)
if err != nil {
    log.WithError(err).Fatal("Failed to load RPC services")
}
```

### **Current Implementation:**
- ✅ **Configuration parsing** - Reads service configs from YAML
- ✅ **Logging** - Logs service loading status
- 🔄 **Dynamic loading** - Placeholder for future implementation
- 📋 **Extensible** - Easy to add new services

## ✅ **Benefits**

### **Configuration-Driven:**
- **Enable/disable services** via YAML config
- **Specify import paths** for each service
- **Clean separation** of service configuration
- **Easy testing** with different service combinations

### **Scalable Architecture:**
- **Add new services** without code changes
- **Modular configuration** for each service
- **Dynamic loading** capability (future)
- **Clear service boundaries**

### **Development Friendly:**
- **Explicit service paths** - Easy to understand
- **Type-safe configuration** - Generated Go structs
- **Logging integration** - Debug service loading
- **Shared config support** - Reuse common settings

## 🎯 **Next Steps**

### **Implement Dynamic Loading:**
1. **Complete RPC loader** with reflection-based loading
2. **Service registration** in gRPC server
3. **Error handling** for missing services
4. **Service discovery** capabilities

### **Add More Services:**
```yaml
rpc_services:
  claw_machine:
    enabled: true
    # ... config
  player:
    enabled: true
    # ... config
  game:                    # NEW SERVICE
    enabled: true
    service_package: "game"
    service_name: "GameService"
    import_path: "github.com/Richard-inter/game/pkg/protocol/game"
    implementation_path: "github.com/Richard-inter/game/internal/service/rpc/game"
```

Your RPC services configuration is now properly structured and ready for dynamic service loading! 🚀
