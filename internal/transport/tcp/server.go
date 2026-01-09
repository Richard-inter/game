package tcp

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Richard-inter/game/internal/config"
	"go.uber.org/zap"
)

const (
	connectionRetryDelay = 100 * time.Millisecond
)

type Server struct {
	config     *config.Config
	logger     *zap.SugaredLogger
	listener   net.Listener
	clients    map[net.Conn]bool
	clientsCh  chan net.Conn
	disconnect chan net.Conn
	messages   chan []byte
}

func NewServer(cfg *config.Config, logger *zap.SugaredLogger) *Server {
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

	s.logger.Infow("Starting TCP server", "address", s.listener.Addr().String())

	// Start connection handler
	go s.handleConnections()

	// Start message broadcaster
	go s.broadcastMessages()

	// Accept connections
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				s.logger.Warnw("Temporary TCP error, retrying...", "error", err)
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

	s.logger.Infow("Shutting down TCP server")

	// Close listener
	if err := s.listener.Close(); err != nil {
		s.logger.Errorw("Failed to close TCP listener", "error", err)
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
	s.logger.Infow("TCP client connected", "client", clientAddr)

	// Set read/write timeouts
	if s.config.TCP.ReadTimeout > 0 {
		if err := conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.TCP.ReadTimeout) * time.Second)); err != nil {
			s.logger.Errorw("Failed to set read deadline", "error", err)
		}
	}
	if s.config.TCP.WriteTimeout > 0 {
		if err := conn.SetWriteDeadline(time.Now().Add(time.Duration(s.config.TCP.WriteTimeout) * time.Second)); err != nil {
			s.logger.Errorw("Failed to set write deadline", "error", err)
		}
	}

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Bytes()

		s.logger.Debugw("Received TCP message", "client", clientAddr, "message", string(message))

		// Echo message back (placeholder)
		if _, err := conn.Write(message); err != nil {
			s.logger.Errorw("Failed to send TCP message", "error", err)
			break
		}

		// Reset read deadline
		if s.config.TCP.ReadTimeout > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(s.config.TCP.ReadTimeout) * time.Second))
		}
	}

	if err := scanner.Err(); err != nil {
		s.logger.Errorw("TCP client error", "error", err, "client", clientAddr)
	}

	s.logger.Infow("TCP client disconnected", "client", clientAddr)
}

func (s *Server) broadcastMessages() {
	for message := range s.messages {
		for client := range s.clients {
			if _, err := client.Write(message); err != nil {
				s.logger.Errorw("Failed to broadcast TCP message", "error", err)
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
		s.logger.Warnw("TCP broadcast channel is full")
	}
}
