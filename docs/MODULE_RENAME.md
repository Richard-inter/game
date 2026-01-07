# Module Rename Guide

## 🔄 **Changed from `github.com/1nterdigital/game` to `github.com/Richard-inter/game`**

### **✅ Files Updated**

#### **1. Go Module**
- `go.mod` - Updated module path

#### **2. Service Main Files**
- `cmd/game-service/main.go` - Updated internal imports
- `cmd/api-service/main.go` - Updated internal imports  
- `cmd/server/main.go` - Updated internal imports
- `cmd/tcp-service/main.go` - Updated internal imports
- `cmd/websocket-service/main.go` - Updated internal imports
- `cmd/server/main_test.go` - Updated internal imports

#### **3. Transport Layer**
- `internal/transport/manager.go` - Updated internal imports
- `internal/transport/http/server.go` - Updated internal imports
- `internal/transport/grpc/server.go` - Updated internal imports
- `internal/transport/tcp/server.go` - Updated internal imports
- `internal/transport/websocket/server.go` - Updated internal imports

#### **4. Service Layer**
- `internal/service/game_grpc_service.go` - Updated protocol imports
- `internal/service/player_grpc_service.go` - Updated protocol imports
- `internal/service/game_service.go` - Updated domain imports

#### **5. Repository Layer**
- `internal/repository/game_repository.go` - Updated domain imports
- `internal/repository/player_repository.go` - Updated domain imports

#### **6. Protocol Files**
- `pkg/protocol/game.proto` - Updated go_package option
- `pkg/protocol/clawMachine/clawMachine.proto` - Updated go_package option
- `pkg/protocol/player/player.proto` - Updated go_package option

#### **7. Generated Files**
- `pkg/protocol/game.pb.go` - Updated package references
- `pkg/protocol/clawMachine/clawMachine.pb.go` - Updated import paths
- `pkg/protocol/player/player.pb.go` - Updated package references

#### **8. Configuration Files**
- `Makefile` - Updated init command
- `.golangci.yml` - Updated linting rules

### **🔄 What Changed**

#### **Before:**
```go
module github.com/1nterdigital/game

import "github.com/1nterdigital/game/internal/config"
import "github.com/1nterdigital/game/pkg/protocol"
```

#### **After:**
```go
module github.com/Richard-inter/game

import "github.com/Richard-inter/game/internal/config"
import "github.com/Richard-inter/game/pkg/protocol"
```

### **🔧 Commands Used**

```bash
# 1. Updated go.mod
module github.com/Richard-inter/game

# 2. Updated all import statements
# Used replace_all to change all occurrences

# 3. Regenerated protobuf files
make proto

# 4. Cleaned up dependencies
go mod tidy
```

### **📋 Verification**

#### **Check if all imports are updated:**
```bash
# Should return no results
grep -r "github.com/1nterdigital/game" .

# Should show new imports
grep -r "github.com/Richard-inter/game" .
```

#### **Test compilation:**
```bash
go build ./...
go test ./...
```

### **🚀 Next Steps**

1. **Push to GitHub:**
   ```bash
   git add .
   git commit -m "Rename module from 1nterdigital/game to Richard-inter/game"
   git push origin main
   ```

2. **Update any external references:**
   - CI/CD pipelines
   - Documentation
   - README files
   - Other repositories that import this module

3. **Test everything:**
   ```bash
   make test
   make build
   make lint
   ```

### **✅ Benefits**

- **Clean namespace** - Uses your GitHub username
- **Better organization** - Clear ownership
- **Professional appearance** - Consistent with your brand
- **Future-proof** - Ready for open source development

All import names have been successfully updated! 🎉
