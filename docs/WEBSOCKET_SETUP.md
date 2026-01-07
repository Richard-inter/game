# WebSocket Setup Guide

## 🎯 **WebSocket Architecture**

This project uses a **clean, dedicated WebSocket service**:

### **Standalone WebSocket Service** (The Only Way)
- **Port**: 8081
- **Path**: `/ws`
- **URL**: `ws://localhost:8081/ws`
- **Service**: Independent WebSocket service with full implementation
- **Purpose**: Real-time bidirectional communication for games

## 🚀 **How to Start WebSocket Service**

### **Method 1: VS Code Launch Configuration**
```bash
# Start infrastructure first
make docker-up-infra

# Launch WebSocket service
# In VS Code: Run & Debug → "Launch WebSocket Service" → F5
```

### **Method 2: Command Line**
```bash
# Start infrastructure
make docker-up-infra

# Start WebSocket service
make run-websocket
```

### **Method 3: Launch All Services**
```bash
# In VS Code: Run & Debug → "Launch All Services" → F5
# This starts: Game, API, WebSocket, and TCP services
```

## 🔧 **Testing WebSocket Connections**

### **Test with WebSocket Client**
```bash
# Using websocat (install with: brew install websocat)
websocat ws://localhost:8081/ws

# Send a message and it will echo back
Hello WebSocket!
# Response: Hello WebSocket!
```

### **Test with Browser JavaScript**
```javascript
// Connect to WebSocket service
const ws = new WebSocket('ws://localhost:8081/ws');

ws.onopen = function(event) {
    console.log('Connected to WebSocket service');
    ws.send('Hello WebSocket!');
};

ws.onmessage = function(event) {
    console.log('Received:', event.data);
};

ws.onclose = function(event) {
    console.log('Disconnected from WebSocket');
};
```

### **Test with curl (Health Check)**
```bash
# Check WebSocket service health
curl http://localhost:8081/health

# Expected response
{
    "status": "healthy",
    "service": "websocket-service"
}
```

## 📋 **Service Status Check**

### **Check if WebSocket is Running**
```bash
# Check port 8081
lsof -i :8081

# Use service monitor script
node scripts/service-monitor.js

# Expected output
✅ WebSocket Service (WebSocket) - Port 8081 - PID: 12345
```

## 🔄 **WebSocket Message Flow**

### **Current Implementation (Echo Server)**
```
Client → WebSocket Server → Echo back to Client
```

1. **Client connects** to `ws://localhost:8081/ws`
2. **Server upgrades** HTTP connection to WebSocket
3. **Client sends message**: `"Hello WebSocket"`
4. **Server receives** and logs the message
5. **Server echoes** back: `"Hello WebSocket"`
6. **Client receives** echoed message

### **Future Implementation (Game Logic)**
```
Client → WebSocket Server → Game Logic → Broadcast to All Clients
```

## 🛠️ **Configuration**

### **WebSocket Service Configuration**
```yaml
# config/websocket-service.yaml
service:
  name: "websocket-service"
  host: "0.0.0.0"
  port: 8081
  path: "/ws"

websocket:
  read_buffer_size: 1024
  write_buffer_size: 1024
  ping_interval: 54
  pong_wait: 60
  write_wait: 10
```

### **Shared Configuration**
```yaml
# config/shared.yaml (inherited by websocket-service)
logging:
  level: "info"
  format: "json"

redis:
  host: "localhost"
  port: 6379
```

## 🔍 **Troubleshooting**

### **Common Issues**

#### **1. WebSocket Connection Refused**
```bash
# Check if service is running
lsof -i :8081

# Start WebSocket service
make run-websocket
```

#### **2. "Upgrade Required" Error**
```bash
# Make sure you're connecting to WebSocket endpoint
# Correct: ws://localhost:8081/ws
# Wrong:  http://localhost:8081/ws
```

#### **3. CORS Issues**
The WebSocket server allows all origins for development:
```go
CheckOrigin: func(_ *http.Request) bool {
    return true // Allow all origins for now
}
```

#### **4. Connection Drops**
Check the logs for WebSocket errors:
```bash
# If running via VS Code, check debug console
# If running via command line, check terminal output
```

### **Debug Commands**
```bash
# Test connection with telnet
telnet localhost 8081

# View WebSocket handshake
curl -i -N -H "Connection: Upgrade" \
     -H "Upgrade: websocket" \
     -H "Sec-WebSocket-Key: test" \
     -H "Sec-WebSocket-Version: 13" \
     http://localhost:8081/ws
```

## 🎮 **Usage Examples**

### **Simple Chat Test**
```javascript
const ws = new WebSocket('ws://localhost:8081/ws');

ws.onopen = () => {
    console.log('Connected!');
    ws.send('Hello from client!');
};

ws.onmessage = (event) => {
    console.log('Echo received:', event.data);
};
```

### **Real-time Game Updates**
```javascript
const gameWs = new WebSocket('ws://localhost:8081/ws');

gameWs.onopen = () => {
    // Send game action
    gameWs.send(JSON.stringify({
        type: 'game_action',
        action: 'move_piece',
        player_id: 'player-123',
        game_id: 'game-456'
    }));
};

gameWs.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Game update:', data);
};
```

## 📊 **Monitoring**

### **WebSocket Metrics**
- **Active Connections**: Number of connected clients
- **Messages Sent/Received**: Message throughput
- **Connection Errors**: Failed connections/disconnections
- **Message Latency**: Time to process messages

### **Health Endpoints**
```bash
# WebSocket service health
curl http://localhost:8081/health

# Expected response
{
    "status": "healthy",
    "service": "websocket-service"
}
```

## 🚀 **Service Architecture**

### **Clean Separation of Concerns**
```
┌─────────────────┐    ┌─────────────────┐
│   HTTP API      │    │  WebSocket      │
│   (Port 8080)   │    │  (Port 8081)    │
│                 │    │                 │
│ • REST endpoints│    │ • Real-time     │
│ • CRUD operations│ │ • Bidirectional  │
│ • Health checks │    │ • Game events   │
└─────────────────┘    └─────────────────┘
```

### **Why This Architecture?**
✅ **Single Responsibility** - Each service has one clear purpose  
✅ **Independent Scaling** - WebSocket can be scaled separately  
✅ **Clean Code** - No mixing of HTTP and WebSocket logic  
✅ **Better Testing** - Each service can be tested independently  
✅ **Easier Maintenance** - Clear boundaries between services  

## 🎯 **Next Steps**

1. **Test Connection**: Use websocat or browser to verify WebSocket works
2. **Implement Game Logic**: Replace echo server with actual game functionality
3. **Add Authentication**: Secure WebSocket connections
4. **Add Rooms**: Support multiple game rooms/channels
5. **Add Persistence**: Store game state in database
6. **Add Monitoring**: Track WebSocket metrics and health

Your WebSocket service is now **clean, focused, and ready for real-time gaming!** 🎉
