# Linting Guide

## Overview

This project uses golangci-lint with a comprehensive configuration copied from aka-im-bot and adapted for the game project structure.

## 🚀 **Getting Started**

### **Install golangci-lint**
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Or use brew
brew install golangci-lint
```

### **Run Linting**
```bash
# Check for issues
make lint

# Auto-fix issues where possible
make lint-fix

# Format code
make fmt
```

## 📋 **Linting Configuration**

### **Enabled Linters**
- **bodyclose** - Checks whether HTTP response body is closed successfully
- **copyloopvar** - Detects loop variables that are copied
- **depguard** - Checks package imports
- **dogsled** - Checks assignments with too many blank identifiers
- **dupl** - Tool for code clone detection
- **errcheck** - Errcheck is a program for checking for unchecked errors
- **errorlint** - errorlint is a linter for Go code that finds coding errors
- **funlen** - Tool for detection of long functions
- **gocheckcompilerdirectives** - linter that checks problematic compiler directives
- **goconst** - Finds repeated strings that could be replaced by a constant
- **gocritic** - The most opinionated Go source code linter
- **gocyclo** - Computes and checks the cyclomatic complexity of functions
- **godox** - Tool for detection of FIXME, TODO and other comment keywords
- **mnd** - Magic number detector
- **goprintffuncname** - Checks that printf-like functions are named with `f` suffix
- **gosec** - Inspects source code for security problems
- **govet** - Vet examines Go source code and reports suspicious constructs
- **intrange** - Finds places where `for range` could be simplified
- **ineffassign** - Detects ineffectual assignments
- **lll** - Long line linter
- **misspell** - Finds commonly misspelled English words
- **nakedret** - Finds naked returns in functions greater than a specified length
- **noctx** - Finds sending http request without context.Context
- **nolintlint** - Reports ill-formed or insufficient nolint directives
- **revive** - Fast, configurable, extensible, flexible, and beautiful linter for Go
- **testifylint** - Checks usage of github.com/stretchr/testify
- **unconvert** - Remove unnecessary type conversions
- **unparam** - Reports unused function parameters
- **unused** - Checks Go code for unused constants, variables, functions and types
- **whitespace** - Tool for detection of leading and trailing whitespace

### **Key Settings**

#### **DepGuard Rules**
- **Logrus**: Only allowed in `pkg/logger` package
- **pkg/errors**: Should use standard library errors
- **instana/testify**: Use stretchr/testify instead

#### **Code Quality**
- **Function length**: Max 50 statements
- **Cyclomatic complexity**: Max 15
- **Line length**: 140 characters
- **Magic numbers**: Detected (0,1,2,3 ignored)

#### **Exclusions**
- **Generated files**: Protocol buffers and FlatBuffers
- **Test files**: `*_test.go`
- **Database patterns**: Repository CRUD operations
- **Service patterns**: Business logic methods
- **Transport patterns**: Server initialization

## 🛠️ **Common Issues & Fixes**

### **1. Unused Parameters**
```go
// ❌ Bad
func setupMiddleware(engine *gin.Engine, cfg *config.ServiceConfig, log *logrus.Logger) {
    // cfg is unused
}

// ✅ Good
func setupMiddleware(engine *gin.Engine, _ *config.ServiceConfig, log *logrus.Logger) {
    // Use underscore for unused params
}
```

### **2. Magic Numbers**
```go
// ❌ Bad
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

// ✅ Good
const shutdownTimeout = 5 * time.Second
ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
```

### **3. HTTP Context Usage**
```go
// ❌ Bad
resp, err := http.Get("http://example.com")

// ✅ Good
req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com", nil)
resp, err := http.DefaultClient.Do(req)
```

### **4. Error Handling**
```go
// ❌ Bad
data, _ := json.Marshal(payload)

// ✅ Good
data, err := json.Marshal(payload)
if err != nil {
    return nil, fmt.Errorf("failed to marshal payload: %w", err)
}
```

### **5. Import Organization**
```go
// ❌ Bad
import (
    "fmt"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/1nterdigital/game/internal/config"
)

// ✅ Good (auto-fixed by goimports)
import (
    "fmt"
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/1nterdigital/game/internal/config"
)
```

## 🎯 **Current Linting Issues**

### **High Priority**
1. **Unused parameters** - Fix with underscore or remove
2. **Magic numbers** - Create constants
3. **HTTP context** - Add context to requests
4. **Error handling** - Handle all errors

### **Medium Priority**
1. **Import organization** - Auto-fix with `make lint-fix`
2. **Code formatting** - Auto-fix with `make fmt`
3. **Function length** - Break down long functions

### **Low Priority**
1. **Code duplication** - Acceptable in some patterns
2. **Security warnings** - Review and fix if needed

## 🔄 **Development Workflow**

### **Before Commit**
```bash
# 1. Format code
make fmt

# 2. Auto-fix linting issues
make lint-fix

# 3. Check remaining issues
make lint

# 4. Run tests
make test
```

### **VS Code Integration**
Add to `.vscode/settings.json`:
```json
{
    "go.lintOnSave": "file",
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
    ]
}
```

### **Pre-commit Hook**
```bash
#!/bin/sh
# .git/hooks/pre-commit
make fmt
make lint-fix
make lint
make test
```

## 📊 **Linting Results Summary**

### **Current Status**: 49 issues found
- **depguard**: 3 issues (package imports)
- **errcheck**: 5 issues (unchecked errors)
- **errorlint**: 1 issue (error handling)
- **gocritic**: 6 issues (code quality)
- **goimports**: 3 issues (import organization)
- **gosec**: 2 issues (security)
- **mnd**: 18 issues (magic numbers)
- **noctx**: 3 issues (HTTP context)
- **revive**: 8 issues (code style)

### **Priority Fix Order**
1. **revive** (unused parameters) - Quick fixes
2. **goimports** - Auto-fixable
3. **mnd** (magic numbers) - Create constants
4. **noctx** (HTTP context) - Add context
5. **errcheck** (error handling) - Handle errors

## 🚀 **Quick Fix Commands**

```bash
# Auto-fix imports and formatting
make lint-fix && make fmt

# Fix specific issues
golangci-lint run --fix --disable-all -E goimports,E revive

# Check specific linters
golangci-lint run --disable-all -E revive,mnd,noctx
```

## 📝 **Best Practices**

1. **Run linting early** in development
2. **Fix issues incrementally** - don't try to fix everything at once
3. **Use meaningful constants** instead of magic numbers
4. **Handle all errors** properly
5. **Use context** for HTTP requests
6. **Remove unused code** - parameters, imports, variables
7. **Keep functions small** and focused
8. **Follow Go conventions** for naming and structure

This linting configuration ensures consistent, high-quality code across the game microservices project!
