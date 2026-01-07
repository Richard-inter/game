package tcp

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"github.com/sirupsen/logrus"
)

const (
	connectionRetryDelay = 100 * time.Millisecond
)

type Server struct {
	config     *config.Config
	logger     *logrus.Logger
	listener   net.Listener
	clients    map[net.Conn]bool
	clientsCh  chan net.Conn
	disconnect chan net.Conn
	messages   chan []byte
}

func NewServer(cfg *config.Config, logger *logrus.Logger) *Server {
	return &Server{
		config:     cfg,
		logger:     logger,
		clients:    make(map[net.Conn]bool),
		clientsCh:  make(chan net.Conn),
		disconnect: make(chan net.Conn),
		messages:   make(chan []byte),
	}
}

func (s *Server) Start() error {
	var err error

	// Create listener
	lc := net.ListenConfig{}
	s.listener, err = lc.Listen(context.Background(), "tcp", fmt.Sprintf("%s:%d", s.config.TCP.Host, s.config.TCP.Port))
	if err != nil {
		return fmt.Errorf("failed to listen on TCP port: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"address": s.listener.Addr().String(),
	}).Info("Starting TCP server")

	// Start connection handler
	go s.handleConnections()

	// Start message broadcaster
	go s.broadcastMessages()

	// Accept connections
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				s.logger.WithError(err).Warn("Temporary TCP error, retrying...")
				time.Sleep(connectionRetryDelay)
				continue
			}
			return fmt.Errorf("TCP accept error: %w", err)
		}

		s.clientsCh <- conn
	}
}

func (s *Server) Shutdown(_ context.Context) error {
	if s.listener == nil {
		return nil
	}

	s.logger.Info("Shutting down TCP server")

	// Close listener
	if err := s.listener.Close(); err != nil {
		s.logger.WithError(err).Error("Failed to close TCP listener")
	}

	// Close all client connections
	for client := range s.clients {
		client.Close()
		delete(s.clients, client)
	}

	return nil
}

func (s *Server) handleConnections() {
	for {
		select {
		case conn := <-s.clientsCh:
			s.clients[conn] = true
			go s.handleClient(conn)
		case conn := <-s.disconnect:
			delete(s.clients, conn)
			conn.Close()
		}
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		s.disconnect <- conn
	}()

	clientAddr := conn.RemoteAddr().String()
	s.logger.WithField("client", clientAddr).Info("TCP client connected")

	// Set read/write timeouts
	if s.config.TCP.ReadTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.TCP.ReadTimeout) * time.Second)); err != nil {
			s.logger.WithError(err).Error("Failed to set read deadline")
		}
	}
	if s.config.TCP.WriteTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(time.Duration(s.config.TCP.WriteTimeout) * time.Second)); err != nil {
			s.logger.WithError(err).Error("Failed to set write deadline")
		}
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()

		s.logger.WithFields(logrus.Fields{
			"client":  clientAddr,
			"message": string(message),
		}).Debug("Received TCP message")

		// Echo message back (placeholder)
		if _, err := conn.Write(message); err != nil {
			s.logger.WithError(err).Error("Failed to send TCP message")
			break
		}

		// Reset read deadline
		if s.config.TCP.ReadTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.TCP.ReadTimeout) * time.Second))
		}
	}

	if err := scanner.Err(); err != nil {
		s.logger.WithError(err).WithField("client", clientAddr).Error("TCP client error")
	}

	s.logger.WithField("client", clientAddr).Info("TCP client disconnected")
}

func (s *Server) broadcastMessages() {
	for message := range s.messages {
		for client := range s.clients {
			if _, err := client.Write(message); err != nil {
				s.logger.WithError(err).Error("Failed to broadcast TCP message")
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

// Broadcast message to all connected clients
func (s *Server) Broadcast(message []byte) {
	select {
	case s.messages <- message:
	default:
		s.logger.Warn("TCP broadcast channel is full")
	}
}
