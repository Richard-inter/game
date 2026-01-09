package main

import (
	"os"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	httptransport "github.com/Richard-inter/game/internal/transport/http"
	"github.com/Richard-inter/game/pkg/logger"
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

	log.Infow("Starting API Service",
		"version", Version,
		"build_time", BuildTime,
		"go_version", GoVersion,
		"service", "api-service",
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/api-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// Create gRPC client manager with service discovery
	var etcdEndpoints []string
	var playerAddr string
	var clawmachineAddr string

	if cfg.Discovery.Enabled && len(cfg.Discovery.Etcd.Endpoints) > 0 {
		etcdEndpoints = cfg.Discovery.Etcd.Endpoints
		log.Infow("Using etcd endpoints from config", "endpoints", etcdEndpoints)
	} else {
		// Service discovery disabled - use direct connections
		log.Infow("Service discovery disabled, using direct gRPC connections")
		etcdEndpoints = []string{}         // Empty to disable service discovery
		playerAddr = "localhost:9094"      // Player service direct address
		clawmachineAddr = "localhost:9091" // Clawmachine service direct address
	}

	grpcClientManager, err := grpc.NewClientManager(&grpc.ClientManagerConfig{
		EtcdEndpoints:   etcdEndpoints,
		PlayerAddr:      playerAddr,
		ClawmachineAddr: clawmachineAddr,
	})
	if err != nil {
		log.Fatalw("Failed to create gRPC client manager", "error", err)
	}
	defer grpcClientManager.Close()

	// Create HTTP server
	server := httptransport.NewServer(cfg, log, grpcClientManager)

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalw("Failed to start API service", "error", err)
	}
}
