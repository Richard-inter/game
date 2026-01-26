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
	"github.com/Richard-inter/game/internal/repository"
	c "github.com/Richard-inter/game/internal/service/rpc/whackAMole_runtime"
	"github.com/Richard-inter/game/pkg/logger"
	pb "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket"
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

	log.Infow("Starting WhackAMole Runtime Service",
		"version", Version,
		"buildTime", BuildTime,
		"goVersion", GoVersion,
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/rpc-whackamole-runtime-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

	// Debug: Log database configuration
	log.Infow("Database configuration loaded",
		"whackAMoleDB", cfg.WhackAMoleDatabase,
		"host", cfg.WhackAMoleDatabase.Host,
		"port", cfg.WhackAMoleDatabase.Port,
		"user", cfg.WhackAMoleDatabase.User,
		"name", cfg.WhackAMoleDatabase.Name,
	)

	// Initialize database
	database, err := db.InitWhackAMoleDB(cfg)
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

	// Initialize and register WhackAMole runtime service
	whackAMoleRepo := repository.NewWhackAMoleRepository(database)

	// Initialize Redis client
	redisClient := cache.NewRedisClient(cfg.GetRedisAddr(), cfg.GetRedisPassword())

	runtimeService := c.NewWhackAMoleWebsocketService(whackAMoleRepo, redisClient, cfg.StreamConsumer.StreamKey)
	pb.RegisterWhackAMoleRuntimeServiceServer(s, runtimeService)

	// Enable reflection for development
	reflection.Register(s)

	log.Infow("WhackAMole Runtime gRPC server starting", "address", lis.Addr().String())

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

	log.Info("Shutting down WhackAMole Runtime service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info("WhackAMole Runtime service stopped gracefully")
	case <-ctx.Done():
		log.Info("Shutdown timeout, forcing stop...")
		s.Stop()
	}
}
