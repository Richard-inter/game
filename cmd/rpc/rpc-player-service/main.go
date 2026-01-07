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

	"github.com/Richard-inter/game/internal/config"
	p "github.com/Richard-inter/game/internal/service/rpc/player"
	"github.com/Richard-inter/game/pkg/logger"
	player "github.com/Richard-inter/game/pkg/protocol/player"
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
		"version":   Version,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
	}).Info("Starting Player Service")

	// Load service-specific configuration
	configFile := os.Getenv("CONFIG_PATH")
	if configFile == "" {
		configFile = "config/rpc-player-service.yaml" // fallback
	}

	cfg, err := config.LoadServiceConfigFromPath(configFile)
	if err != nil {
		log.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize gRPC server
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", cfg.GetGRPCAddr())
	if err != nil {
		log.WithError(err).Fatal("Failed to listen")
	}

	s := grpc.NewServer()

	// Initialize and register player service
	playerService := p.NewPlayerGRPCService()
	player.RegisterPlayerServiceServer(s, playerService)

	// Enable reflection for development
	reflection.Register(s)

	log.WithFields(logrus.Fields{
		"address": lis.Addr().String(),
	}).Info("Player gRPC server starting")

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

	log.Info("Shutting down Player service...")

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
		log.Info("Player service stopped gracefully")
	case <-ctx.Done():
		log.Info("Player service shutdown timeout")
		s.Stop()
	}
}
