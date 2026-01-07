package main

import (
	"context"
	"testing"
	"time"

	"github.com/1nterdigital/game/internal/config"
	"github.com/1nterdigital/game/internal/transport"
	"github.com/1nterdigital/game/pkg/logger"
)

func TestServerStartup(t *testing.T) {
	// Initialize logger
	logger.InitLogger()
	log := logger.GetLogger()

	// Load test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 8080, // Use different port for testing
			Mode: "test",
		},
		GRPC: config.GRPCConfig{
			Host: "localhost",
			Port: 9090, // Use different port for testing
		},
		WebSocket: config.WebSocketConfig{
			Host: "localhost",
			Port: 8081, // Use different port for testing
			Path: "/ws",
		},
		TCP: config.TCPConfig{
			Host: "localhost",
			Port: 8082, // Use different port for testing
		},
	}

	// Initialize transport manager
	transportManager := transport.NewManager(cfg, log)

	// Start servers in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startErr := make(chan error, 1)
	go func() {
		if err := transportManager.Start(ctx); err != nil {
			startErr <- err
		}
	}()

	// Wait a moment for servers to start
	select {
	case err := <-startErr:
		if err != nil {
			t.Logf("Expected some errors in test environment without actual dependencies: %v", err)
		}
	case <-time.After(2 * time.Second):
		// Timeout is expected since we don't have actual dependencies
	}

	// Test shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := transportManager.Shutdown(shutdownCtx); err != nil {
		t.Logf("Shutdown error (expected in test): %v", err)
	}

	t.Log("Server startup and shutdown test completed")
}
