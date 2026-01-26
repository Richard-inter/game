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

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/db"
	"github.com/Richard-inter/game/internal/repository"
	whackAMole "github.com/Richard-inter/game/internal/service/rpc/whackAMole"
	"github.com/Richard-inter/game/pkg/logger"
	whackAMolepb "github.com/Richard-inter/game/pkg/protocol/whackAMole"
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

	log.Infow("Starting WhackAMole Service",
		"version", Version,
		"buildTime", BuildTime,
		"goVersion", GoVersion,
	)

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/rpc-whackamole-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.Fatalw("Failed to load configuration", "error", err)
	}

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

	// Initialize and register whack a mole service
	whackAMoleRepo := repository.NewWhackAMoleRepository(database)
	whackAMoleService := whackAMole.NewWhackAMoleGRPCService(whackAMoleRepo)
	whackAMolepb.RegisterWhackAMoleServiceServer(s, whackAMoleService)

	// Enable reflection for development
	reflection.Register(s)

	log.Infow("WhackAMole gRPC server starting", "address", lis.Addr().String())

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

	log.Infow("Shutting down WhackAMole service...")

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
		log.Infow("WhackAMole service stopped gracefully")
	case <-ctx.Done():
		log.Infow("WhackAMole service shutdown timeout")
		s.Stop()
	}
}
