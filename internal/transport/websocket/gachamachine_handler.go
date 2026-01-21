package websocket

import (
	"context"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/transport/grpc"
	runtimepb "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket/gachaMachine"
)

type GachaMachineWebSocketHandler struct {
	logger       *zap.SugaredLogger
	gachaRuntime *grpc.GachaMachineRuntimeClient

	handlers map[fbs.MessageType]messageHandler
}

func NewGachaMachineWebSocketHandler(logger *zap.SugaredLogger, grpcManager *grpc.ClientManager) (*GachaMachineWebSocketHandler, error) {
	runtimeClient, err := grpcManager.GetGachaMachineRuntimeClient()
	if err != nil {
		return nil, err
	}

	h := &GachaMachineWebSocketHandler{
		logger:       logger,
		gachaRuntime: runtimeClient,
		handlers:     make(map[fbs.MessageType]messageHandler),
	}

	h.handlers[fbs.MessageTypeGetPullResultWsReq] = h.handleGetPullResult
	h.handlers[fbs.MessageTypeGetPlayerInfoWsReq] = h.handleGetPlayerInfo
	h.handlers[fbs.MessageTypeGetMachineInfoWsReq] = h.handleGetMachineInfo

	return h, nil
}

func (h *GachaMachineWebSocketHandler) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	h.logger.Infow("GachaMachine WebSocket client connected")

	for {
		// Read message
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Errorw("GachaMachine WebSocket error", "error", err)
			} else {
				h.logger.Infow("GachaMachine WebSocket connection closed", "error", err)
			}
			return
		}

		h.logger.Debugw("Received GachaMachine WebSocket message", "message_type", messageType, "message", string(message))

		// Handle message
		response, err := h.handleMessage(message)
		if err != nil {
			h.logger.Errorw("Error handling GachaMachine message", "error", err)
			h.sendError(conn, "Failed to process message")
			continue
		}

		// Send response
		if err := conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
			h.logger.Errorw("Error sending GachaMachine response", "error", err)
			return
		}
	}
}

func (h *GachaMachineWebSocketHandler) handleMessage(data []byte) ([]byte, error) {
	envelope := fbs.GetRootAsEnvelope(data, 0)
	msgType := envelope.Type()
	payload := envelope.PayloadBytes()

	handler, ok := h.handlers[msgType]
	if !ok {
		h.logger.Errorw("Unknown GachaMachine message type", "type", msgType)
		return h.buildErrorResp(400, "Unknown message type"), nil
	}

	if len(payload) == 0 {
		h.logger.Errorw("Empty payload for GachaMachine", "type", msgType)
		return h.buildErrorResp(400, "Empty payload"), nil
	}

	return handler(context.Background(), payload)
}

func (h *GachaMachineWebSocketHandler) handleGetPullResult(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.gachaRuntime.GetPullResultWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("GetPullResultWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *GachaMachineWebSocketHandler) handleGetPlayerInfo(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.gachaRuntime.GetPlayerInfoWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("GetPlayerInfoWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *GachaMachineWebSocketHandler) handleGetMachineInfo(
	ctx context.Context,
	payload []byte,
) ([]byte, error) {
	resp, err := h.gachaRuntime.GetMachineInfoWs(ctx, &runtimepb.RuntimeRequest{
		Payload: payload,
	})
	if err != nil {
		h.logger.Errorw("GetMachineInfoWs failed", "error", err)
		return h.buildErrorResp(500, err.Error()), nil
	}

	return resp.Payload, nil
}

func (h *GachaMachineWebSocketHandler) buildErrorResp(code int32, message string) []byte {
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

func (h *GachaMachineWebSocketHandler) sendError(conn *websocket.Conn, message string) {
	response := h.buildErrorResp(500, message)

	err := conn.WriteMessage(websocket.BinaryMessage, response)
	if err != nil {
		h.logger.Errorw("Failed to send GachaMachine error message", "error", err)
	}
}
