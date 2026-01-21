package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	wshandler "github.com/Richard-inter/game/internal/transport/websocket"
	"github.com/Richard-inter/game/pkg/logger"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	shutdownTimeout = 5 * time.Second
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetSugar()

	log.Infow("Starting WebSocket Service",
		"version", Version,
		"build_time", BuildTime,
		"go_version", GoVersion,
		"service", "websocket-service",
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/websocket-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// WebSocket upgrader
	upgrader := websocket.Upgrader{
		ReadBufferSize:  cfg.WebSocket.ReadBufferSize,
		WriteBufferSize: cfg.WebSocket.WriteBufferSize,
		CheckOrigin: func(_ *http.Request) bool {
			return true // Allow all origins for now
		},
	}

	// Create HTTP server
	mux := http.NewServeMux()

	// ClawMachine WebSocket endpoint
	mux.HandleFunc("/clawmachine", func(w http.ResponseWriter, r *http.Request) {
		handleClawMachineWebSocket(upgrader, w, r, log)
	})

	// GachaMachine WebSocket endpoint
	mux.HandleFunc("/gachamachine", func(w http.ResponseWriter, r *http.Request) {
		handleGachaMachineWebSocket(upgrader, w, r, log)
	})

	// Add health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"websocket-service"}`)
	})

	server := &http.Server{
		Addr:         cfg.GetServiceAddr(),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infow("WebSocket Service starting", "address", server.Addr, "path", cfg.WebSocket.Path)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Failed to start WebSocket service", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutting down WebSocket Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorw("WebSocket service shutdown error", "error", err)
	}

	log.Infow("WebSocket Service stopped")
}

func handleClawMachineWebSocket(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, log *zap.SugaredLogger) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorw("Failed to upgrade ClawMachine WebSocket connection", "error", err)
		return
	}

	// Load service configurations dynamically
	serviceConfigs, err := config.LoadMultipleServiceConfigs([]string{"player", "clawmachine", "clawmachine_runtime"})
	if err != nil {
		log.Errorw("Failed to load service configs", "error", err)
		return
	}

	// Initialize gRPC client manager with dynamic addresses
	grpcConfig := &grpc.ClientManagerConfig{
		EtcdEndpoints:   []string{}, // Empty to disable etcd
		PlayerAddr:      fmt.Sprintf("%s:%d", serviceConfigs["player"].Service.Host, serviceConfigs["player"].Service.Port),
		ClawmachineAddr: fmt.Sprintf("%s:%d", serviceConfigs["clawmachine"].Service.Host, serviceConfigs["clawmachine"].Service.Port),
		RuntimeAddr:     fmt.Sprintf("%s:%d", serviceConfigs["clawmachine_runtime"].Service.Host, serviceConfigs["clawmachine_runtime"].Service.Port),
	}

	log.Infow("Creating ClawMachine gRPC client manager with dynamic addresses",
		"player", grpcConfig.PlayerAddr,
		"clawmachine", grpcConfig.ClawmachineAddr,
		"clawmachine_runtime", grpcConfig.RuntimeAddr)

	grpcManager, err := grpc.NewClientManager(grpcConfig)
	if err != nil {
		log.Errorw("Failed to create gRPC client manager", "error", err)
		return
	}

	log.Infow("ClawMachine gRPC client manager created successfully")

	// Create WebSocket handler
	handler, err := wshandler.NewClawMachineWebSocketHandler(log, grpcManager)
	if err != nil {
		log.Errorw("Failed to create ClawMachine WebSocket handler", "error", err)
		return
	}

	// Handle the connection
	handler.HandleConnection(conn)
}

func handleGachaMachineWebSocket(upgrader websocket.Upgrader, w http.ResponseWriter, r *http.Request, log *zap.SugaredLogger) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorw("Failed to upgrade GachaMachine WebSocket connection", "error", err)
		return
	}

	// Load service configurations dynamically
	serviceConfigs, err := config.LoadMultipleServiceConfigs([]string{"gachamachine", "gachamachine_runtime"})
	if err != nil {
		log.Errorw("Failed to load GachaMachine service configs", "error", err)
		return
	}

	// Initialize gRPC client manager with dynamic addresses
	grpcConfig := &grpc.ClientManagerConfig{
		EtcdEndpoints:    []string{}, // Empty to disable etcd
		GachaMachineAddr: fmt.Sprintf("%s:%d", serviceConfigs["gachamachine"].Service.Host, serviceConfigs["gachamachine"].Service.Port),
		GachaRuntimeAddr: fmt.Sprintf("%s:%d", serviceConfigs["gachamachine_runtime"].Service.Host, serviceConfigs["gachamachine_runtime"].Service.Port),
	}

	log.Infow("Creating GachaMachine gRPC client manager with dynamic addresses",
		"gachamachine", grpcConfig.GachaMachineAddr,
		"gachamachine_runtime", grpcConfig.GachaRuntimeAddr)

	grpcManager, err := grpc.NewClientManager(grpcConfig)
	if err != nil {
		log.Errorw("Failed to create GachaMachine gRPC client manager", "error", err)
		return
	}

	log.Infow("GachaMachine gRPC client manager created successfully")

	// Create GachaMachine WebSocket handler
	handler, err := wshandler.NewGachaMachineWebSocketHandler(log, grpcManager)
	if err != nil {
		log.Errorw("Failed to create GachaMachine WebSocket handler", "error", err)
		return
	}

	// Handle the connection
	handler.HandleConnection(conn)
}
