# MySQL Docker Setup Guide

## 🎯 **Simple MySQL-Only Docker Setup**

### **Files Created:**
- `docker-compose.mysql.yml` - MySQL only configuration
- Updated `Makefile` - Added infrastructure commands

## 🐳 **Docker Compose (MySQL Only)**

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    container_name: game-mysql
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=game
      - MYSQL_USER=gameuser
      - MYSQL_PASSWORD=gamepass
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./data/mysql/init:/docker-entrypoint-initdb.d
    networks:
      - game-network
    restart: unless-stopped

volumes:
  mysql_data:

networks:
  game-network:
    driver: bridge
```

## 🔧 **Makefile Commands**

### **Infrastructure Commands:**
```bash
# Start MySQL only
make docker-up-infra

# Stop MySQL only  
make docker-down-infra

# View MySQL logs
make docker-logs-infra

# Clean up MySQL resources
make docker-clean
```

### **Full Commands (for reference):**
```bash
# Start all services (including game services)
make docker-up

# Stop all services
make docker-down

# View all logs
make docker-logs
```

## 🚀 **Usage**

### **1. Start MySQL:**
```bash
make docker-up-infra
```

### **2. Start Game Services (VS Code):**
```bash
# Start infrastructure first
make docker-up-infra

# Then launch services in VS Code:
# - Launch Game Service (Port 9090)
# - Launch API Service (Port 8080)  
# - Launch WebSocket Service (Port 8081)
# - Launch TCP Service (Port 8082)
```

### **3. Connect to MySQL:**
```bash
# From your Go services
# Host: mysql (Docker network)
# Port: 3306
# Database: game
# User: gameuser
# Password: gamepass

# From external tools
# Host: localhost
# Port: 3306
```

### **4. Verify MySQL:**
```bash
# Check if MySQL is running
docker-compose -f docker-compose.mysql.yml ps

# Connect to MySQL container
docker exec -it game-mysql mysql -u root -p

# Check logs
make docker-logs-infra
```

## 📋 **MySQL Configuration**

### **Environment Variables:**
- `MYSQL_ROOT_PASSWORD=rootpassword` - Root password
- `MYSQL_DATABASE=game` - Database name
- `MYSQL_USER=gameuser` - Application user
- `MYSQL_PASSWORD=gamepass` - Application password

### **Connection Details:**
- **Host**: `mysql` (Docker network) or `localhost:3306` (host)
- **Port**: `3306`
- **Database**: `game`
- **User**: `gameuser`
- **Password**: `gamepass`

## 🔍 **Troubleshooting**

### **Common Issues:**

#### **1. Port Already in Use:**
```bash
# Check what's using port 3306
lsof -i :3306

# Kill process if needed
sudo kill -9 <PID>
```

#### **2. Connection Refused:**
```bash
# Check if MySQL container is running
docker-compose -f docker-compose.mysql.yml ps

# Check logs for errors
make docker-logs-infra
```

#### **3. Database Not Created:**
```bash
# Check init scripts
ls -la data/mysql/init/

# Manually create database if needed
docker exec -it game-mysql mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS game;"
```

## 🎉 **Benefits**

✅ **Lightweight** - Only runs MySQL, no extra services  
✅ **Fast startup** - Quick to start and stop  
✅ **Resource efficient** - Minimal memory usage  
✅ **Simple debugging** - Easy to troubleshoot  
✅ **Focused development** - Perfect for microservice development  

Your MySQL setup is now simplified and ready for development! 🚀
