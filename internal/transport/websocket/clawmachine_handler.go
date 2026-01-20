package websocket

import (
	"context"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	runtimepb "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket/clawMachine"
)

type messageHandler func(ctx context.Context, payload []byte) ([]byte, error)

type ClawMachineWebSocketHandler struct {
	logger            *zap.SugaredLogger
	playerClient      *grpc.PlayerClient
	clawmachineClient *grpc.ClawMachineClient
	wsClient          *grpc.ClawMachineRuntimeClient

	handlers map[fbs.MessageType]messageHandler
}

func NewClawMachineWebSocketHandler(logger *zap.SugaredLogger, grpcManager *grpc.ClientManager) (*ClawMachineWebSocketHandler, error) {
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

	h := &ClawMachineWebSocketHandler{
		logger:            logger,
		playerClient:      playerClient,
		clawmachineClient: clawmachineClient,
		wsClient:          runtimeClient,
		handlers:          make(map[fbs.MessageType]messageHandler),
	}

	h.handlers[fbs.MessageTypeStartClawGameReq] = h.handleStartClawGame
	h.handlers[fbs.MessageTypeGetPlayerInfoWsReq] = h.handleGetPlayerInfo
	h.handlers[fbs.MessageTypeAddTouchedItemRecordReq] = h.handleAddTouchedItemRecord
	h.handlers[fbs.MessageTypeSpawnItemReq] = h.handleSpawnItem

	return h, nil
}

func (h *ClawMachineWebSocketHandler) HandleConnection(conn *websocket.Conn) {
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
		if err := conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
			h.logger.Errorw("Error sending response", "error", err)
			return
		}
	}
}

func (h *ClawMachineWebSocketHandler) handleMessage(data []byte) ([]byte, error) {
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

func (h *ClawMachineWebSocketHandler) handleStartClawGame(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.wsClient.StartClawGameWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("StartClawGameWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *ClawMachineWebSocketHandler) handleGetPlayerInfo(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.wsClient.GetPlayerSnapshotWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("GetPlayerSnapshotWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *ClawMachineWebSocketHandler) handleAddTouchedItemRecord(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.wsClient.AddTouchedItemRecordWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("AddTouchedItemRecordWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *ClawMachineWebSocketHandler) handleSpawnItem(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.wsClient.SpawnItemWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("SpawnItemWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *ClawMachineWebSocketHandler) buildErrorResp(code int32, message string) []byte {
	builder := flatbuffers.NewBuilder(128)

	msgOffset := builder.CreateString(message)

	fbs.ErrorRespStart(builder)
	fbs.ErrorRespAddCode(builder, code)
	fbs.ErrorRespAddMessage(builder, msgOffset)
	errorResp := fbs.ErrorRespEnd(builder)

	builder.Finish(errorResp)
	errorBytes := builder.FinishedBytes()

	// Wrap error response in Envelope
	envBuilder := flatbuffers.NewBuilder(256)
	payloadOffset := envBuilder.CreateByteVector(errorBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeErrorResp)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	return envBuilder.FinishedBytes()
}

func (h *ClawMachineWebSocketHandler) sendError(conn *websocket.Conn, message string) {
	response := h.buildErrorResp(500, message)

	err := conn.WriteMessage(websocket.BinaryMessage, response)
	if err != nil {
		h.logger.Errorw("Failed to send error message", "error", err)
	}
}
