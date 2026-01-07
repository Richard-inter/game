package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/1nterdigital/game/internal/config"
	"github.com/1nterdigital/game/internal/transport/grpc"
	httptransport "github.com/1nterdigital/game/internal/transport/http"
	"github.com/1nterdigital/game/internal/transport/tcp"
	"github.com/1nterdigital/game/internal/transport/websocket"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	config     *config.Config
	logger     *logrus.Logger
	httpServer *httptransport.Server
	grpcServer *grpc.Server
	wsServer   *websocket.Server
	tcpServer  *tcp.Server
	wg         sync.WaitGroup
}

func NewManager(cfg *config.Config, logger *logrus.Logger) *Manager {
	return &Manager{
		config: cfg,
		logger: logger,
	}
}

func (m *Manager) Start(_ context.Context) error {
	// Start HTTP Server
	m.httpServer = httptransport.NewServer(m.config, m.logger)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		if err := m.httpServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			m.logger.WithError(err).Error("HTTP server failed to start")
		}
	}()

	// Start gRPC Server
	m.grpcServer = grpc.NewServer(m.config, m.logger)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		if err := m.grpcServer.Start(); err != nil {
			m.logger.WithError(err).Error("gRPC server failed to start")
		}
	}()

	// Start WebSocket Server
	m.wsServer = websocket.NewServer(m.config, m.logger)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		if err := m.wsServer.Start(); err != nil {
			m.logger.WithError(err).Error("WebSocket server failed to start")
		}
	}()

	// Start TCP Server
	m.tcpServer = tcp.NewServer(m.config, m.logger)
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		if err := m.tcpServer.Start(); err != nil {
			m.logger.WithError(err).Error("TCP server failed to start")
		}
	}()

	m.logger.Info("All transport servers started successfully")
	return nil
}

func (m *Manager) Shutdown(ctx context.Context) error {
	m.logger.Info("Shutting down transport manager...")

	var errs []error

	// Shutdown HTTP Server
	if m.httpServer != nil {
		if err := m.httpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("HTTP server shutdown error: %w", err))
		}
	}

	// Shutdown gRPC Server
	if m.grpcServer != nil {
		if err := m.grpcServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("gRPC server shutdown error: %w", err))
		}
	}

	// Shutdown WebSocket Server
	if m.wsServer != nil {
		if err := m.wsServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("WebSocket server shutdown error: %w", err))
		}
	}

	// Shutdown TCP Server
	if m.tcpServer != nil {
		if err := m.tcpServer.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("TCP server shutdown error: %w", err))
		}
	}

	// Wait for all servers to shutdown
	m.wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	m.logger.Info("Transport manager shutdown completed")
	return nil
}
