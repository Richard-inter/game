package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/db"
	"github.com/Richard-inter/game/internal/discovery"
	"github.com/Richard-inter/game/internal/registry"
	"github.com/Richard-inter/game/internal/repository"
	g "github.com/Richard-inter/game/internal/service/rpc/gachaMachine"
	"github.com/Richard-inter/game/pkg/logger"
	gachaMachine "github.com/Richard-inter/game/pkg/protocol/gachaMachine"
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

	log.Infow("Starting GachaMachine Service",
		"version", Version,
		"buildTime", BuildTime,
		"goVersion", GoVersion,
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/rpc-gachamachine-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// Initialize database
	database, err := db.InitGachaMachineDB(cfg)
	if err != nil {
		log.Fatalw("Failed to initialize database", "error", err)
	}

	// Initialize gRPC server
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", cfg.GetGRPCAddr())
	if err != nil {
		log.Fatalw("Failed to listen", "error", err)
	}

	s := grpc.NewServer()

	// Initialize and register gacha machine service
	gachaMachineRepo := repository.NewGachaMachineRepository(database)

	// Initialize Redis client
	redisClient := cache.NewRedisClient(cfg.GetRedisAddr(), cfg.GetRedisPassword())

	gachaMachineService := g.NewGachaMachineGRPCService(gachaMachineRepo, redisClient)
	gachaMachine.RegisterGachaMachineServiceServer(s, gachaMachineService)

	// Enable reflection for development
	reflection.Register(s)

	log.Infow("GachaMachine gRPC server starting", "address", lis.Addr().String())

	// Register service with etcd
	if cfg.Discovery.Enabled {
		etcdDiscovery, err := discovery.NewEtcdDiscovery(cfg.Discovery.Etcd.Endpoints)
		if err != nil {
			log.Warnw("Failed to connect to etcd, service discovery disabled", "error", err)
		} else {
			serviceRegistry := registry.NewServiceRegistry(etcdDiscovery, log)
			err = serviceRegistry.RegisterService("gachamachine-service", cfg.Service.Host, cfg.Service.Port)
			if err != nil {
				log.Errorw("Failed to register service with etcd", "error", err)
			}
			defer etcdDiscovery.Close()
		}
	}

	// Start server in a goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalw("Failed to serve", "error", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Infow("Shutting down GachaMachine service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Infow("GachaMachine service stopped gracefully")
	case <-ctx.Done():
		log.Infow("GachaMachine service shutdown timeout")
		s.Stop()
	}
}
