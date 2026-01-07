# VS Code Launch Configuration Guide

## 🎯 **Updated Launch Configurations**

Updated `.vscode/launch.json` to support the new RPC services structure with proper naming and organization.

## 🚀 **New Launch Configurations**

### **Individual Service Launch:**

#### **1. RPC ClawMachine Service**
```json
{
    "name": "Launch RPC ClawMachine Service",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}/cmd/rpc/rpc-clawmachine-service",
    "env": {
        "CONFIG_PATH": "${workspaceFolder}/config/rpc-clawmachine-service.yaml"
    },
    "cwd": "${workspaceFolder}",
    "args": [],
    "showLog": true
}
```

#### **2. RPC Player Service**
```json
{
    "name": "Launch RPC Player Service",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}/cmd/rpc/rpc-player-service",
    "env": {
        "CONFIG_PATH": "${workspaceFolder}/config/rpc-player-service.yaml"
    },
    "cwd": "${workspaceFolder}",
    "args": [],
    "showLog": true
}
```

### **Compound Launch Configurations:**

#### **1. Launch All Services**
Launches all 6 services simultaneously:
- RPC ClawMachine Service (Port 9091)
- RPC Player Service (Port 9092)
- Game Service (Port 9090)
- API Service (Port 8080)
- WebSocket Service (Port 8081)
- TCP Service (Port 8082)

#### **2. Launch RPC Services Only**
Launches only the RPC services:
- RPC ClawMachine Service (Port 9091)
- RPC Player Service (Port 9092)

## 🔧 **VS Code Debugging Features**

### **Environment Variables:**
Each service automatically sets the correct `CONFIG_PATH`:
- **RPC ClawMachine**: `config/rpc-clawmachine-service.yaml`
- **RPC Player**: `config/rpc-player-service.yaml`
- **Game Service**: `config/game-service.yaml`
- **API Service**: `config/api-service.yaml`
- **WebSocket**: `config/websocket-service.yaml`
- **TCP**: `config/tcp-service.yaml`

### **Debug Features:**
- **`showLog: true`** - Shows service logs in debug console
- **`mode: "auto"`** - Automatic debug mode detection
- **`cwd: "${workspaceFolder}"`** - Sets working directory
- **Individual configs** - Each service has isolated configuration

## 🎯 **Usage in VS Code**

### **Debug Panel Access:**
1. **Open Debug Panel** → `Ctrl+Shift+D` (or `Run → Start Debugging`)
2. **Select Configuration** → Choose from dropdown:
   - `Launch RPC ClawMachine Service`
   - `Launch RPC Player Service`
   - `Launch All Services`
   - `Launch RPC Services Only`
3. **Start Debugging** → Press `F5` or click green play button
4. **Set Breakpoints** → Click in code to set breakpoints
5. **View Logs** → Debug console shows service logs

### **Service-Specific Debugging:**

#### **RPC ClawMachine Service:**
- **Entry Point**: `cmd/rpc/rpc-clawmachine-service/main.go`
- **Config**: `config/rpc-clawmachine-service.yaml`
- **Port**: 9091
- **Protocol**: gRPC with reflection enabled

#### **RPC Player Service:**
- **Entry Point**: `cmd/rpc/rpc-player-service/main.go`
- **Config**: `config/rpc-player-service.yaml`
- **Port**: 9092
- **Protocol**: gRPC with reflection enabled

## 🔍 **Debugging Workflow**

### **1. Start Infrastructure:**
```bash
make docker-up-infra
```

### **2. Debug Individual Services:**
- **Open `cmd/rpc/rpc-clawmachine-service/main.go`**
- **Set breakpoints** in service code
- **Select "Launch RPC ClawMachine Service"**
- **Press F5** to start debugging

### **3. Test Services:**
```bash
# Test RPC ClawMachine service
grpcurl -plaintext localhost:9091 describe

# Test RPC Player service
grpcurl -plaintext localhost:9092 describe

# Check service health
curl http://localhost:8080/health
```

## ✅ **Benefits**

### **Development Efficiency:**
- **One-click debugging** - Launch services directly from VS Code
- **Integrated debugging** - Breakpoints, step-through, variable inspection
- **Log visibility** - Service logs appear in debug console
- **Configuration management** - Auto-sets correct config paths

### **Service Isolation:**
- **Independent debugging** - Debug one service at a time
- **Separate configs** - Each service uses its own configuration
- **Clear boundaries** - No shared state between services
- **Focused development** - Work on specific service functionality

### **Team Collaboration:**
- **Multiple developers** - Each can debug different services
- **Consistent setup** - Same launch configuration for team
- **Easy onboarding** - Clear service structure for new developers
- **Parallel development** - Work on multiple services simultaneously

## 🎮 **Debugging Tips**

### **gRPC Service Debugging:**
1. **Use reflection** - Services have reflection enabled for testing
2. **Check logs** - Service startup logs appear in debug console
3. **Test endpoints** - Use `grpcurl` to test gRPC services
4. **Network issues** - Check port conflicts if services don't start

### **Breakpoints Strategy:**
- **Main function** - Set at service entry point
- **Service methods** - Set in RPC service implementations
- **Error handling** - Set in error handling code paths
- **Configuration loading** - Set at config parsing

Your VS Code launch configuration is now optimized for the new RPC services structure! 🚀
