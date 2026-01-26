package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket/whackAMole"
	flatbuffers "github.com/google/flatbuffers/go"
)

func main() {
	// Connect to the WhackAMole runtime service
	conn, err := grpc.NewClient("localhost:9098", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewWhackAMoleRuntimeServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("🎮 Testing WhackAMole Runtime Service")
	fmt.Println("=====================================")

	// Test 1: GetMoleWeight
	fmt.Println("\n📊 Test 1: GetMoleWeight Request")
	testGetMoleWeight(ctx, client)

	// Test 2: GetLeaderboard (without player_id)
	fmt.Println("\n🏆 Test 2: GetLeaderboard Request (no player_id)")
	testGetLeaderboard(ctx, client, 0)

	// Test 3: GetLeaderboard (with player_id)
	fmt.Println("\n🏆 Test 3: GetLeaderboard Request (with player_id=123)")
	testGetLeaderboard(ctx, client, 123)

	fmt.Println("\n✅ All tests completed!")
}

func testGetMoleWeight(ctx context.Context, client pb.WhackAMoleRuntimeServiceClient) {
	// Create request using helper function
	reqBytes := createGetMoleWeightRequest()

	// Create request
	req := &pb.RuntimeRequest{
		Payload: reqBytes,
	}

	// Send request
	resp, err := client.GetMoleWeight(ctx, req)
	if err != nil {
		log.Printf("❌ GetMoleWeight failed: %v", err)
		return
	}

	// Parse response
	respEnvelope := fbs.GetRootAsEnvelope(resp.Payload, 0)
	if respEnvelope == nil {
		log.Printf("❌ Invalid response envelope")
		return
	}

	if respEnvelope.Type() != fbs.MessageTypeGetMoleWeightResp {
		log.Printf("❌ Wrong response type: %d", respEnvelope.Type())
		return
	}

	// Parse GetMoleWeightResp
	respData := fbs.GetRootAsGetMoleWeightResp(respEnvelope.PayloadBytes(), 0)
	if respData == nil {
		log.Printf("❌ Invalid GetMoleWeightResp")
		return
	}

	// Print mole weight configs
	moleVector := respData.MoleLength()
	fmt.Printf("✅ GetMoleWeight successful! Found %d mole weight configs:\n", moleVector)

	for i := 0; i < moleVector; i++ {
		var mole fbs.MoleWeight
		if respData.Mole(&mole, i) {
			fmt.Printf("  - Mole %d: Type='%s', Weight=%d\n", i+1, mole.MoleType(), mole.Weight())
		}
	}
}

func testGetLeaderboard(ctx context.Context, client pb.WhackAMoleRuntimeServiceClient, playerID int64) {
	// Create request using helper function
	reqBytes := createGetLeaderboardRequest(playerID)

	// Create request
	req := &pb.RuntimeRequest{
		Payload: reqBytes,
	}

	// Send request
	resp, err := client.GetLeaderboard(ctx, req)
	if err != nil {
		log.Printf("❌ GetLeaderboard failed: %v", err)
		return
	}

	// Parse response
	respEnvelope := fbs.GetRootAsEnvelope(resp.Payload, 0)
	if respEnvelope == nil {
		log.Printf("❌ Invalid response envelope")
		return
	}

	if respEnvelope.Type() != fbs.MessageTypeGetLeaderboardResp {
		log.Printf("❌ Wrong response type: %d", respEnvelope.Type())
		return
	}

	// Parse GetLeaderboardResp
	respData := fbs.GetRootAsGetLeaderboardResp(respEnvelope.PayloadBytes(), 0)
	if respData == nil {
		log.Printf("❌ Invalid GetLeaderboardResp")
		return
	}

	// Print leaderboard data
	playersVector := respData.TopPlayerLength()
	yourRank := respData.YourRank()
	yourScore := respData.YourScore()

	fmt.Printf("✅ GetLeaderboard successful! Top %d players:\n", playersVector)

	if playerID > 0 {
		fmt.Printf("  Your Rank: %d, Your Score: %d\n", yourRank, yourScore)
	}

	for i := 0; i < playersVector; i++ {
		var player fbs.LeaderboardPlayer
		if respData.TopPlayer(&player, i) {
			fmt.Printf("  - Rank %d: PlayerID=%d, Username='%s', Score=%d\n",
				player.Rank(), player.PlayerId(), player.Username(), player.Score())
		}
	}
}

// createGetMoleWeightRequest builds a GetMoleWeightReq Envelope
func createGetMoleWeightRequest() []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build GetMoleWeightReq table
	fbs.GetMoleWeightReqStart(builder)
	req := fbs.GetMoleWeightReqEnd(builder)
	builder.Finish(req)

	// Get the serialized request bytes
	reqBytes := builder.FinishedBytes()

	// Create new builder for envelope
	envBuilder := flatbuffers.NewBuilder(512)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetMoleWeightReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	env := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(env)

	return envBuilder.FinishedBytes()
}

// createGetLeaderboardRequest builds a GetLeaderboardReq Envelope
func createGetLeaderboardRequest(playerID int64) []byte {
	builder := flatbuffers.NewBuilder(1024)

	// Build GetLeaderboardReq table
	fbs.GetLeaderboardReqStart(builder)
	fbs.GetLeaderboardReqAddPlayerId(builder, playerID)
	req := fbs.GetLeaderboardReqEnd(builder)
	builder.Finish(req)

	// Get the serialized request bytes
	reqBytes := builder.FinishedBytes()

	// Create new builder for envelope
	envBuilder := flatbuffers.NewBuilder(512)
	payloadOffset := envBuilder.CreateByteVector(reqBytes)

	fbs.EnvelopeStart(envBuilder)
	fbs.EnvelopeAddType(envBuilder, fbs.MessageTypeGetLeaderboardReq)
	fbs.EnvelopeAddPayload(envBuilder, payloadOffset)
	env := fbs.EnvelopeEnd(envBuilder)
	envBuilder.Finish(env)

	return envBuilder.FinishedBytes()
}
