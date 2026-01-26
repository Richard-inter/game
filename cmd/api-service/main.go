package main

import (
	"fmt"
	"os"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	httptransport "github.com/Richard-inter/game/internal/transport/http"
	"github.com/Richard-inter/game/pkg/logger"
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
	var gachaMachineAddr string

	if cfg.Discovery.Enabled && len(cfg.Discovery.Etcd.Endpoints) > 0 {
		etcdEndpoints = cfg.Discovery.Etcd.Endpoints
		log.Infow("Using etcd endpoints from config", "endpoints", etcdEndpoints)
	} else {
		// Service discovery disabled - load required service configs
		log.Infow("Service discovery disabled, loading service configurations")
		etcdEndpoints = []string{} // Empty to disable service discovery

		serviceConfigs, err := config.LoadMultipleServiceConfigs([]string{"player", "clawmachine", "gachamachine"})
		if err != nil {
			log.Fatalw("Failed to load service configs", "error", err)
		}

		playerAddr = fmt.Sprintf("%s:%d", serviceConfigs["player"].Service.Host, serviceConfigs["player"].Service.Port)
		clawmachineAddr = fmt.Sprintf("%s:%d", serviceConfigs["clawmachine"].Service.Host, serviceConfigs["clawmachine"].Service.Port)
		gachaMachineAddr = fmt.Sprintf("%s:%d", serviceConfigs["gachamachine"].Service.Host, serviceConfigs["gachamachine"].Service.Port)

		log.Infow("Using gRPC service addresses from service configs",
			"player", playerAddr,
			"clawmachine", clawmachineAddr,
			"gachamachine", gachaMachineAddr)
	}

	grpcClientManager, err := grpc.NewClientManager(&grpc.ClientManagerConfig{
		EtcdEndpoints:         etcdEndpoints,
		PlayerAddr:            playerAddr,
		ClawmachineAddr:       clawmachineAddr,
		GachaMachineAddr:      gachaMachineAddr,
		WhackAMoleAddr:        "localhost:9097", // WhackAMole service address
		WhackAMoleRuntimeAddr: "localhost:9098", // WhackAMole runtime service address
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
