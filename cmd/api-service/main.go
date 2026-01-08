package main

import (
	"os"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	httptransport "github.com/Richard-inter/game/internal/transport/http"
	"github.com/Richard-inter/game/pkg/logger"
	"github.com/sirupsen/logrus"
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
	log := logger.GetLogger()

	log.WithFields(logrus.Fields{
		"version":    Version,
		"build_time": BuildTime,
		"go_version": GoVersion,
		"service":    "api-service",
	}).Info("Starting API Service")

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/api-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Create gRPC client manager with service discovery
	var etcdEndpoints []string
	var playerAddr string
	var clawmachineAddr string

	if cfg.Discovery.Enabled && len(cfg.Discovery.Etcd.Endpoints) > 0 {
		etcdEndpoints = cfg.Discovery.Etcd.Endpoints
		log.WithField("endpoints", etcdEndpoints).Info("Using etcd endpoints from config")
	} else {
		// Service discovery disabled - use direct connections
		log.Info("Service discovery disabled, using direct gRPC connections")
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
		log.WithError(err).Fatal("Failed to create gRPC client manager")
	}
	defer grpcClientManager.Close()

	// Create HTTP server
	server := httptransport.NewServer(cfg, log, grpcClientManager)

	// Start server
	if err := server.Start(); err != nil {
		log.WithError(err).Fatal("Failed to start API service")
	}
}
