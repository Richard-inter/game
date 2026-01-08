package websocket

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/Richard-inter/game/internal/transport/grpc"
	player "github.com/Richard-inter/game/pkg/protocol/player"
)

type WebSocketHandler struct {
	logger       *logrus.Logger
	playerClient *grpc.PlayerClient
}

func NewWebSocketHandler(logger *logrus.Logger, grpcManager *grpc.ClientManager) (*WebSocketHandler, error) {
	playerClient, err := grpcManager.GetPlayerClient()
	if err != nil {
		return nil, err
	}

	return &WebSocketHandler{
		logger:       logger,
		playerClient: playerClient,
	}, nil
}

func (h *WebSocketHandler) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	h.logger.Info("WebSocket client connected")

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.WithError(err).Error("WebSocket error")
			} else {
				h.logger.WithError(err).Info("WebSocket connection closed")
			}
			return
		}

		h.logger.WithFields(logrus.Fields{
			"message_type": messageType,
			"message":      string(message),
		}).Debug("Received WebSocket message")

		// Handle message
		response, err := h.handleMessage(message)
		if err != nil {
			h.logger.WithError(err).Error("Error handling message")
			h.sendError(conn, "Failed to process message")
			continue
		}

		// Send response
		if err := conn.WriteMessage(messageType, response); err != nil {
			h.logger.WithError(err).Error("Error sending response")
			return
		}
	}
}

func (h *WebSocketHandler) handleMessage(data []byte) ([]byte, error) {
	// For now, just echo the message back
	// TODO: Implement proper message parsing and handling
	h.logger.Info("Processing WebSocket message")
	return data, nil
}

func (h *WebSocketHandler) handlePlayerRequest(playerID int64) ([]byte, error) {
	// Example: Get player info via gRPC
	resp, err := h.playerClient.GetPlayerInfo(context.Background(), &player.GetPlayerInfoReq{
		PlayerID: playerID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get player info: %w", err)
	}

	// Convert response to bytes (simplified for now)
	return []byte(fmt.Sprintf("Player info: %+v", resp.Player)), nil
}

func (h *WebSocketHandler) sendError(conn *websocket.Conn, message string) {
	response := []byte(fmt.Sprintf(`{"error": "%s"}`, message))

	err := conn.WriteMessage(websocket.TextMessage, response)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send error message")
	}
}
