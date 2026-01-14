package websocket

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	runtimepb "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket/clawMachine"
	player "github.com/Richard-inter/game/pkg/protocol/player"
)

type WebSocketHandler struct {
	logger            *zap.SugaredLogger
	playerClient      *grpc.PlayerClient
	clawmachineClient *grpc.ClawMachineClient
	wsClient          *grpc.ClawMachineRuntimeClient
}

func NewWebSocketHandler(logger *zap.SugaredLogger, grpcManager *grpc.ClientManager) (*WebSocketHandler, error) {
	playerClient, err := grpcManager.GetPlayerClient()
	if err != nil {
		return nil, err
	}

	clawmachineClient, err := grpcManager.GetClawMachineClient()
	if err != nil {
		return nil, err
	}

	runtimeClient, err := grpcManager.GetClawMachineRuntimeClient()
	if err != nil {
		return nil, err
	}

	return &WebSocketHandler{
		logger:            logger,
		playerClient:      playerClient,
		clawmachineClient: clawmachineClient,
		wsClient:          runtimeClient,
	}, nil
}

func (h *WebSocketHandler) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	h.logger.Infow("WebSocket client connected")

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Errorw("WebSocket error", "error", err)
			} else {
				h.logger.Infow("WebSocket connection closed", "error", err)
			}
			return
		}

		h.logger.Debugw("Received WebSocket message", "message_type", messageType, "message", string(message))

		// Handle message
		response, err := h.handleMessage(message)
		if err != nil {
			h.logger.Errorw("Error handling message", "error", err)
			h.sendError(conn, "Failed to process message")
			continue
		}

		// Send response
		if err := conn.WriteMessage(messageType, response); err != nil {
			h.logger.Errorw("Error sending response", "error", err)
			return
		}
	}
}

func (h *WebSocketHandler) handleMessage(data []byte) ([]byte, error) {
	h.logger.Infow("Processing WebSocket message")

	// Parse the payload as FlatBuffer Envelope
	envelope := fbs.GetRootAsEnvelope(data, 0)

	// Check the message type
	msgType := envelope.Type()
	h.logger.Infow("Message type", "type", msgType)

	switch msgType {
	case fbs.MessageTypeStartClawGameReq:
		// Get the payload bytes and parse as StartClawGameReq
		payloadBytes := envelope.PayloadBytes()
		if len(payloadBytes) > 0 {
			startReq := fbs.GetRootAsStartClawGameReq(payloadBytes, 0)
			playerID := startReq.PlayerId()
			machineID := startReq.MachineId()

			h.logger.Infow("StartClawGame request",
				"playerID", playerID,
				"machineID", machineID)

			// Call websocket service and return raw response
			h.logger.Infow("Calling StartClawGameWs", "payload_length", len(payloadBytes))
			wsReq := &runtimepb.RuntimeRequest{
				Payload: payloadBytes,
			}

			resp, err := h.wsClient.StartClawGameWs(context.Background(), wsReq)
			if err != nil {
				h.logger.Errorw("Failed to call StartClawGameWs", "error", err)
				return []byte(`{"error":"Failed to start game"}`), nil
			}

			h.logger.Infow("Received response from service", "response_length", len(resp.Payload))
			h.logger.Infow("Response hex", "hex", fmt.Sprintf("%x", resp.Payload))

			// Return the raw response payload from the service
			return resp.Payload, nil
		} else {
			h.logger.Errorw("Empty StartClawGame payload")
			return []byte(`{"error":"Empty payload"}`), nil
		}

	default:
		h.logger.Errorw("Unknown message type", "type", msgType)
		return []byte(`{"error":"Unknown message type"}`), nil
	}
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
		h.logger.Errorw("Failed to send error message", "error", err)
	}
}
