package main

import (
	"fmt"
	"log"

	fbs "github.com/Richard-inter/game/pkg/protocol/gachaMachine_Websocket/gachaMachine"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/gorilla/websocket"
)

// change to main to test
func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8081/gachamachine", nil)
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
			case fbs.MessageTypeGetPullResultWsResp:
				handleGetPullResultWsResp(env)
				return // stop loop after receiving response
			case fbs.MessageTypeGetPlayerInfoWsResp:
				handleGetPlayerInfoWsResp(env)
				return // stop loop after receiving response
			case fbs.MessageTypeGetMachineInfoWsResp:
				handleGetMachineInfoWsResp(env)
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
	// msg := createGetPullResultWsRequest(1, 1, 1)
	// msg := createGetPlayerInfoWsRequest(1)
	msg := createGetMachineInfoWsRequest(1)
	if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
		log.Fatal("write error:", err)
	}

	// Wait until response is received
	<-done
	fmt.Println("Response received, client exiting.")
}

// createGetPullResultWsRequest builds a GetPullResultWsReq Envelope
func createGetPullResultWsRequest(playerID, machineID int64, pullCount int) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build GetPullResultWsReq table
	fbs.GetPullResultWsReqStart(builder)
	fbs.GetPullResultWsReqAddPlayerId(builder, playerID)
	fbs.GetPullResultWsReqAddMachineId(builder, machineID)
	fbs.GetPullResultWsReqAddPullCount(builder, int32(pullCount))
	req := fbs.GetPullResultWsReqEnd(builder)
	builder.Finish(req)

	// Get the serialized request bytes
	reqBytes := builder.FinishedBytes()

	// Create new builder for envelope
	envBuilder := flatbuffers.NewBuilder(512)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetPullResultWsReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	env := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(env)

	return envBuilder.FinishedBytes()
}

// createGetPlayerInfoWsRequest builds a GetPlayerInfoWsReq Envelope
func createGetPlayerInfoWsRequest(playerID int64) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build GetPlayerInfoWsReq table
	fbs.GetPlayerInfoWsReqStart(builder)
	fbs.GetPlayerInfoWsReqAddPlayerId(builder, playerID)
	req := fbs.GetPlayerInfoWsReqEnd(builder)
	builder.Finish(req)

	// Get the serialized request bytes
	reqBytes := builder.FinishedBytes()

	// Create new builder for envelope
	envBuilder := flatbuffers.NewBuilder(512)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetPlayerInfoWsReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	env := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(env)

	return envBuilder.FinishedBytes()
}

// createGetMachineInfoWsRequest builds a GetMachineInfoWsReq Envelope
func createGetMachineInfoWsRequest(machineID int64) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build GetMachineInfoWsReq table
	fbs.GetMachineInfoWsReqStart(builder)
	fbs.GetMachineInfoWsReqAddMachineId(builder, machineID)
	req := fbs.GetMachineInfoWsReqEnd(builder)
	builder.Finish(req)

	// Get the serialized request bytes
	reqBytes := builder.FinishedBytes()

	// Create new builder for envelope
	envBuilder := flatbuffers.NewBuilder(512)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetMachineInfoWsReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	env := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(env)

	return envBuilder.FinishedBytes()
}

// handleGetPullResultWsResp safely reads GetPullResultWsResp from Envelope
func handleGetPullResultWsResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsGetPullResultWsResp(env.PayloadBytes(), 0)

	fmt.Println("=== GetPullResultWsResp ===")
	fmt.Println("Items:")
	for i := 0; i < resp.ItemIdsLength(); i++ {
		fmt.Printf("  Item ID: %d\n", resp.ItemIds(i))
	}
	fmt.Println("==============================")
}

// handleGetPlayerInfoWsResp safely reads GetPlayerInfoWsResp from Envelope
func handleGetPlayerInfoWsResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsGetPlayerInfoWsResp(env.PayloadBytes(), 0)

	fmt.Println("=== GetPlayerInfoWsResp ===")
	fmt.Printf("Player ID: %d\n", resp.PlayerId())
	fmt.Printf("Username: %s\n", string(resp.Username()))
	fmt.Printf("Coin: %d\n", resp.Coin())
	fmt.Printf("Diamond: %d\n", resp.Diamond())
	fmt.Println("============================")
}

// handleGetMachineInfoWsResp safely reads GetMachineInfoWsResp from Envelope
func handleGetMachineInfoWsResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsGetMachineInfoWsResp(env.PayloadBytes(), 0)

	fmt.Println("=== GetMachineInfoWsResp ===")
	fmt.Printf("Machine ID: %d\n", resp.MachineId())
	fmt.Printf("Name: %s\n", string(resp.Name()))
	fmt.Printf("Price: %d\n", resp.Price())
	fmt.Printf("Price Times Ten: %d\n", resp.PriceTimesTen())
	fmt.Printf("Super Rare Pity: %d\n", resp.SuperRarePity())
	fmt.Printf("Ultra Rare Pity: %d\n", resp.UltraRarePity())
	fmt.Println("============================")
}

// handleErrorResp safely reads ErrorResp from Envelope
func handleErrorResp(env *fbs.Envelope) {
	resp := fbs.GetRootAsErrorResp(env.PayloadBytes(), 0)

	fmt.Println("=== ErrorResp ===")
	fmt.Printf("Code: %d\n", resp.Code())
	fmt.Printf("Message: %s\n", string(resp.Message()))
	fmt.Println("==================")
}
