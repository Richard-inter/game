package websocket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/Richard-inter/game/internal/transport/grpc"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config      *config.Config
	logger      *logrus.Logger
	server      *http.Server
	grpcManager *grpc.ClientManager
	upgrader    websocket.Upgrader
	clients     map[*websocket.Conn]bool
}

func NewServer(cfg *config.Config, logger *logrus.Logger, grpcManager *grpc.ClientManager) *Server {
	return &Server{
		config:      cfg,
		logger:      logger,
		grpcManager: grpcManager,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  cfg.WebSocket.ReadBufferSize,
			WriteBufferSize: cfg.WebSocket.WriteBufferSize,
			CheckOrigin: func(_ *http.Request) bool {
				return true // Allow all origins for now
			},
		},
		clients: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) Start() error {
	// Create WebSocket handler
	wsHandler, err := NewWebSocketHandler(s.logger, s.grpcManager)
	if err != nil {
		return fmt.Errorf("failed to create WebSocket handler: %w", err)
	}

	// Create HTTP server for WebSocket
	mux := http.NewServeMux()
	mux.HandleFunc(s.config.WebSocket.Path, func(w http.ResponseWriter, r *http.Request) {
		// Upgrade HTTP connection to WebSocket
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			s.logger.WithError(err).Error("Failed to upgrade WebSocket connection")
			return
		}

		// Handle connection using the handler
		wsHandler.HandleConnection(conn)
	})

	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.WebSocket.Host, s.config.WebSocket.Port),
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	s.logger.WithFields(logrus.Fields{
		"address": s.server.Addr,
		"path":    s.config.WebSocket.Path,
	}).Info("Starting WebSocket server")

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("WebSocket server failed to start: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	s.logger.Info("Shutting down WebSocket server")

	// Close all client connections
	for client := range s.clients {
		client.Close()
		delete(s.clients, client)
	}

	return s.server.Shutdown(ctx)
}

// Broadcast message to all connected clients
func (s *Server) Broadcast(message []byte) {
	for client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			s.logger.WithError(err).Error("Failed to broadcast message")
			client.Close()
			delete(s.clients, client)
		}
	}
}
