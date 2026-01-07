package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/1nterdigital/game/internal/config"
	"github.com/1nterdigital/game/internal/service"
	"github.com/1nterdigital/game/pkg/logger"
	"github.com/1nterdigital/game/pkg/protocol"
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
		"service":    "game-service",
	}).Info("Starting Game Service")

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/game-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize gRPC server
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", cfg.GetServiceAddr())
	if err != nil {
		log.WithError(err).Fatal("Failed to listen")
	}

	s := grpc.NewServer()

	// Initialize services (placeholder implementations)
	gameService := service.NewGameGRPCService()
	playerService := service.NewPlayerGRPCService()

	// Register services
	protocol.RegisterGameServiceServer(s, gameService)
	protocol.RegisterPlayerServiceServer(s, playerService)

	// Enable reflection for development
	reflection.Register(s)

	log.WithFields(logrus.Fields{
		"address": lis.Addr().String(),
	}).Info("Game Service gRPC server starting")

	// Start server in a goroutine
	go func() {
		if err := s.Serve(lis); err != nil {
			log.WithError(err).Fatal("Failed to serve")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Game Service...")

	// Graceful shutdown
	stopped := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Info("Game Service stopped gracefully")
	case <-time.After(shutdownTimeout):
		log.Warn("Game Service shutdown timeout, forcing stop")
		s.Stop()
	}
}
