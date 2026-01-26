package websocket

import (
	"context"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	runtimepb "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket/whackAMole"
)

type WhackAMoleWebSocketHandler struct {
	logger   *zap.SugaredLogger
	wsClient *grpc.WhackAMoleRuntimeClient

	handlers map[fbs.MessageType]messageHandler
}

func NewWhackAMoleWebSocketHandler(logger *zap.SugaredLogger, grpcManager *grpc.ClientManager) (*WhackAMoleWebSocketHandler, error) {
	runtimeClient, err := grpcManager.GetWhackAMoleRuntimeClient()
	if err != nil {
		return nil, err
	}

	h := &WhackAMoleWebSocketHandler{
		logger:   logger,
		wsClient: runtimeClient,
		handlers: make(map[fbs.MessageType]messageHandler),
	}

	h.handlers[fbs.MessageTypeGetMoleWeightReq] = h.handleGetMoleWeight
	h.handlers[fbs.MessageTypeGetLeaderboardReq] = h.handleGetLeaderboard

	return h, nil
}

func (h *WhackAMoleWebSocketHandler) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	h.logger.Infow("WhackAMole WebSocket client connected")

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Errorw("WhackAMole WebSocket error", "error", err)
			} else {
				h.logger.Infow("WhackAMole WebSocket connection closed", "error", err)
			}
			return
		}

		h.logger.Debugw("Received WhackAMole WebSocket message", "message_type", messageType, "message", string(message))

		// Handle message
		response, err := h.handleMessage(message)
		if err != nil {
			h.logger.Errorw("Error handling WhackAMole message", "error", err)
			h.sendError(conn, "Failed to process message")
			continue
		}

		// Send response
		if err := conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
			h.logger.Errorw("Error sending WhackAMole response", "error", err)
			return
		}
	}
}

func (h *WhackAMoleWebSocketHandler) handleMessage(data []byte) ([]byte, error) {
	envelope := fbs.GetRootAsEnvelope(data, 0)
	msgType := envelope.Type()
	payload := envelope.PayloadBytes()

	handler, ok := h.handlers[msgType]
	if !ok {
		h.logger.Errorw("Unknown message type", "type", msgType)
		return h.buildErrorResp(400, "Unknown message type"), nil
	}

	if len(payload) == 0 {
		h.logger.Errorw("Empty payload", "type", msgType)
		return h.buildErrorResp(400, "Empty payload"), nil
	}

	return handler(context.Background(), payload)
}

func (h *WhackAMoleWebSocketHandler) handleGetMoleWeight(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.wsClient.GetMoleWeight(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("GetMoleWeight failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *WhackAMoleWebSocketHandler) handleGetLeaderboard(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.wsClient.GetLeaderboard(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("GetLeaderboard failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *WhackAMoleWebSocketHandler) buildErrorResp(code int, message string) []byte {
	builder := flatbuffers.NewBuilder(0)
	messageOffset := builder.CreateString(message)

	fbs.ErrorRespStart(builder)
	fbs.ErrorRespAddCode(builder, int32(code))
	fbs.ErrorRespAddMessage(builder, messageOffset)
	errorOffset := fbs.ErrorRespEnd(builder)

	// Create envelope
	fbs.EnvelopeStart(builder)
	fbs.EnvelopeAddType(builder, fbs.MessageTypeErrorResp)
	fbs.EnvelopeAddPayload(builder, errorOffset)
	envOffset := fbs.EnvelopeEnd(builder)
	builder.Finish(envOffset)

	return builder.FinishedBytes()
}

func (h *WhackAMoleWebSocketHandler) sendError(conn *websocket.Conn, message string) {
	errorResp := h.buildErrorResp(500, message)
	if err := conn.WriteMessage(websocket.BinaryMessage, errorResp); err != nil {
		h.logger.Errorw("Error sending error response", "error", err)
	}
}
