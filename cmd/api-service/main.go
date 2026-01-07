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

	// Create gRPC client
	grpcClient, err := grpc.NewClient(&grpc.Config{
		PlayerServiceAddr:      "localhost:9094",
		ClawMachineServiceAddr: "localhost:9091",
	})
	if err != nil {
		log.WithError(err).Fatal("Failed to create gRPC client")
	}
	defer grpcClient.Close()

	// Create HTTP server
	server := httptransport.NewServer(cfg, log, grpcClient)

	// Start server
	if err := server.Start(); err != nil {
		log.WithError(err).Fatal("Failed to start API service")
	}
}
