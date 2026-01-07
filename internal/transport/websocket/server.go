package websocket

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config   *config.Config
	logger   *logrus.Logger
	server   *http.Server
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

func NewServer(cfg *config.Config, logger *logrus.Logger) *Server {
	return &Server{
		config: cfg,
		logger: logger,
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
	// Create HTTP server for WebSocket
	mux := http.NewServeMux()
	mux.HandleFunc(s.config.WebSocket.Path, s.handleWebSocket)

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

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.WithError(err).Error("Failed to upgrade WebSocket connection")
		return
	}
	defer conn.Close()

	// Add client to the list
	s.clients[conn] = true
	defer delete(s.clients, conn)

	s.logger.WithFields(logrus.Fields{
		"client_ip":  r.RemoteAddr,
		"user_agent": r.UserAgent(),
	}).Info("WebSocket client connected")

	// Handle WebSocket messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.logger.WithError(err).Error("WebSocket error")
			}
			break
		}

		s.logger.WithFields(logrus.Fields{
			"message_type": messageType,
			"message":      string(message),
		}).Debug("Received WebSocket message")

		// Echo message back (placeholder)
		if err := conn.WriteMessage(messageType, message); err != nil {
			s.logger.WithError(err).Error("Failed to send WebSocket message")
			break
		}
	}

	s.logger.WithField("client_ip", r.RemoteAddr).Info("WebSocket client disconnected")
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
