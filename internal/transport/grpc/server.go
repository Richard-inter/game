package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/1nterdigital/game/internal/config"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Server struct {
	config *config.Config
	logger *logrus.Logger
	server *grpc.Server
}

func NewServer(cfg *config.Config, logger *logrus.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
	}
}

func (s *Server) Start() error {
	// Create listener
	lc := net.ListenConfig{}
	lis, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("%s:%d", s.config.GRPC.Host, s.config.GRPC.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port: %w", err)
	}

	// Create gRPC server
	s.server = grpc.NewServer()

	// Register services (placeholder)
	// pb.RegisterGameServiceServer(s.server, &gameService{})
	// pb.RegisterPlayerServiceServer(s.server, &playerService{})

	s.logger.WithFields(logrus.Fields{
		"address": lis.Addr().String(),
	}).Info("Starting gRPC server")

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("Shutting down gRPC server")

	// Create a channel to signal when shutdown is complete
	done := make(chan struct{})

	// Graceful stop in a goroutine
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	// Wait for graceful stop or timeout
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		// Force stop if timeout is reached
		s.server.Stop()
		return ctx.Err()
	}
}
