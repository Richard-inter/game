package main

import (
	"fmt"
	"log"

	fbs "github.com/Richard-inter/game/pkg/protocol/clawMachine_Websocket/clawMachine"
	flatbuffers "github.com/google/flatbuffers/go"
)

// CreateStartClawGameRequest creates a WebSocket request message for starting a claw game
func CreateStartClawGameRequest(playerID, machineID uint64) []byte {
	// Build FlatBuffer message
	builder := flatbuffers.NewBuilder(1024)

	// Build StartClawGameReq payload
	fbs.StartClawGameReqStart(builder)
	fbs.StartClawGameReqAddPlayerId(builder, playerID)
	fbs.StartClawGameReqAddMachineId(builder, machineID)
	startReqOffset := fbs.StartClawGameReqEnd(builder)

	// Create payload bytes
	builder.Finish(startReqOffset)
	payloadBytes := builder.FinishedBytes()

	// Build Envelope with proper payload
	builder2 := flatbuffers.NewBuilder(1024)
	payloadOffset := builder2.CreateByteVector(payloadBytes)

	fbs.EnvelopeStart(builder2)
	fbs.EnvelopeAddType(builder2, fbs.MessageTypeStartClawGameReq)
	fbs.EnvelopeAddPayload(builder2, payloadOffset)
	envelopeOffset := fbs.EnvelopeEnd(builder2)

	builder2.Finish(envelopeOffset)
	message := builder2.FinishedBytes()

	return message
}

// Example usage function (call this from your main function)
func ExampleCreateRequest() {
	playerID := uint64(12345)
	machineID := uint64(67890)

	message := CreateStartClawGameRequest(playerID, machineID)

	fmt.Printf("WebSocket Request Message:\n")
	fmt.Printf("Player ID: %d\n", playerID)
	fmt.Printf("Machine ID: %d\n", machineID)
	fmt.Printf("Message Type: %d (StartClawGameReq)\n", fbs.MessageTypeStartClawGameReq)
	fmt.Printf("Message Bytes: %v\n", message)
	fmt.Printf("Message Hex: %x\n", message)

	// To send this message to the WebSocket server:
	// 1. Connect to ws://localhost:8081/ws
	// 2. Send the message as binary data: websocket.BinaryMessage
	// 3. The server will respond with JSON format

	log.Println("Request message created successfully!")
}
