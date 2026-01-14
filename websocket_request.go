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
			default:
				log.Println("unknown message type:", env.Type())
			}
		}
	}()

	// Send request
	msg := createStartClawGameRequest(1, 1)
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
