package main

import (
	"fmt"
	"log"

	fbs "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket/clawMachine"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	done := make(chan struct{}) // channel to signal when response is received

	// Start read loop
	go func() {
		defer close(done) // signal when finished
		for {
			msgType, data, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			if msgType != websocket.BinaryMessage {
				log.Println("ignored non-binary message")
				continue
			}

			env := fbs.GetRootAsEnvelope(data, 0)
			switch env.Type() {
			case fbs.MessageTypeStartClawGameResp:
				handleStartClawGameResp(env)
				return // stop loop after receiving response
			case fbs.MessageTypeGetPlayerInfoWsResp:
				handleGetPlayerInfoWsResp(env)
				return // stop loop after receiving response
			case fbs.MessageTypeAddTouchedItemRecordResp:
				handleAddTouchedItemRecordResp(env)
				return // stop loop after receiving response
			case fbs.MessageTypeErrorResp:
				handleErrorResp(env)
				return // stop loop after receiving response
			default:
				log.Printf("unknown message type: %d", env.Type())
				return // stop loop to prevent infinite waiting
			}
		}
	}()

	// Send request
	// msg := createStartClawGameRequest(1, 1)
	// msg := createGetPlayerInfoWsRequest(1)
	msg := createAddTouchedItemRecordRequest(23, 1, true)
	if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		log.Fatal("write error:", err)
	}

	// Wait until response is received
	<-done
	fmt.Println("Response received, client exiting.")
}

// createStartClawGameRequest builds a StartClawGameReq Envelope
func createStartClawGameRequest(playerID, machineID uint64) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build StartClawGameReq table
	fbs.StartClawGameReqStart(builder)
	fbs.StartClawGameReqAddPlayerId(builder, playerID)
	fbs.StartClawGameReqAddMachineId(builder, machineID)
	reqOffset := fbs.StartClawGameReqEnd(builder)
	builder.Finish(reqOffset)
	reqBytes := builder.FinishedBytes()

	// Wrap in Envelope
	envBuilder := flatbuffers.NewBuilder(1024)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeStartClawGameReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	return envBuilder.FinishedBytes()
}

// createGetPlayerInfoWsRequest builds a GetPlayerInfoWsReq Envelope
func createGetPlayerInfoWsRequest(playerID uint64) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build GetPlayerInfoWsReq table
	fbs.GetPlayerInfoWsReqStart(builder)
	fbs.GetPlayerInfoWsReqAddPlayerId(builder, playerID)
	reqOffset := fbs.GetPlayerInfoWsReqEnd(builder)
	builder.Finish(reqOffset)
	reqBytes := builder.FinishedBytes()

	// Wrap in Envelope
	envBuilder := flatbuffers.NewBuilder(1024)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetPlayerInfoWsReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	return envBuilder.FinishedBytes()
}

// createAddTouchedItemRecordRequest builds a AddTouchedItemRecordReq Envelope
func createAddTouchedItemRecordRequest(gameID, itemID uint64, catched bool) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build AddTouchedItemRecordReq table
	fbs.AddTouchedItemRecordReqStart(builder)
	fbs.AddTouchedItemRecordReqAddGameId(builder, gameID)
	fbs.AddTouchedItemRecordReqAddItemId(builder, itemID)
	fbs.AddTouchedItemRecordReqAddCatched(builder, catched)
	reqOffset := fbs.AddTouchedItemRecordReqEnd(builder)
	builder.Finish(reqOffset)
	reqBytes := builder.FinishedBytes()

	// Wrap in Envelope
	envBuilder := flatbuffers.NewBuilder(1024)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeAddTouchedItemRecordReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(envOffset)

	return envBuilder.FinishedBytes()
}

// handleStartClawGameResp safely reads StartClawGameResp from the Envelope
func handleStartClawGameResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsStartClawGameResp(env.PayloadBytes(), 0)

	fmt.Println("=== StartClawGameResp ===")
	fmt.Println("Game ID:", resp.GameId())
	fmt.Println("Results:")

	for i := 0; i < resp.ResultsLength(); i++ {
		var result fbs.ClawResult
		if resp.Results(&result, i) {
			fmt.Printf("  Item %d, catched=%v\n", result.ItemId(), result.Catched())
		}
	}
	fmt.Println("=========================")
}

// handleGetPlayerInfoWsResp safely reads GetPlayerInfoWsResp from the Envelope
func handleGetPlayerInfoWsResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsGetPlayerInfoWsResp(env.PayloadBytes(), 0)

	fmt.Println("=== GetPlayerInfoWsResp ===")
	fmt.Println("Player ID:", resp.PlayerId())
	fmt.Println("Username:", string(resp.Username()))
	fmt.Println("Coin:", resp.Coin())
	fmt.Println("Diamond:", resp.Diamond())
	fmt.Println("==========================")
}

// handleAddTouchedItemRecordResp safely reads AddTouchedItemRecordResp from the Envelope
func handleAddTouchedItemRecordResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsAddTouchedItemRecordResp(env.PayloadBytes(), 0)

	fmt.Println("=== AddTouchedItemRecordResp ===")
	fmt.Println("Game ID:", resp.GameId())
	fmt.Println("Item ID:", resp.ItemId())
	fmt.Println("Catched:", resp.Catched())
	fmt.Println("================================")
}

// handleErrorResp safely reads ErrorResp from the Envelope
func handleErrorResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsErrorResp(env.PayloadBytes(), 0)

	fmt.Println("=== ErrorResp ===")
	fmt.Println("Code:", resp.Code())
	fmt.Println("Message:", string(resp.Message()))
	fmt.Println("==================")
}
