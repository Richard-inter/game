package whackAMole_runtime

import (
	"context"
	"errors"
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
	"go.uber.org/zap"

	"github.com/Richard-inter/game/internal/cache"
	"github.com/Richard-inter/game/internal/repository"
	"github.com/Richard-inter/game/pkg/logger"
	pb "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket"
	fbs "github.com/Richard-inter/game/pkg/protocol/whackAMole_Websocket/whackAMole"
)

type WhackAMoleWebsocketService struct {
	pb.UnimplementedWhackAMoleRuntimeServiceServer
	repo      repository.WhackAMoleRepository
	redis     *cache.RedisClient
	streamKey string
	log       *zap.SugaredLogger
}

func NewWhackAMoleWebsocketService(repo repository.WhackAMoleRepository, redis *cache.RedisClient, streamKey string) *WhackAMoleWebsocketService {
	return &WhackAMoleWebsocketService{
		repo:      repo,
		redis:     redis,
		streamKey: streamKey,
		log:       logger.GetSugar(),
	}
}

func (_ *WhackAMoleWebsocketService) buildEnvelopeResponse(messageType fbs.MessageType, payloadBytes []byte) *pb.RuntimeResponse {
	builder := flatbuffers.NewBuilder(len(payloadBytes) + 256)
	payloadOffset := builder.CreateByteVector(payloadBytes)

	fbs.EnvelopeStart(builder)
	fbs.EnvelopeAddType(builder, messageType)
	fbs.EnvelopeAddPayload(builder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(builder)
	builder.Finish(envOffset)

	return &pb.RuntimeResponse{
		Payload: builder.FinishedBytes(),
	}
}

func (s *WhackAMoleWebsocketService) GetMoleWeight(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	s.log.Info("GetMoleWeight called")

	// Parse the request payload
	reqEnvelope := fbs.GetRootAsEnvelope(req.Payload, 0)
	if reqEnvelope == nil {
		return nil, errors.New("invalid request envelope")
	}

	// Get mole weight configs from repository
	configs, err := s.repo.GetMoleWeightConfig(ctx, 0) // 0 means get all
	if err != nil {
		s.log.Errorw("Failed to get mole weight configs", "error", err)
		return nil, fmt.Errorf("failed to get mole weight configs: %w", err)
	}

	// Build flatbuffers response
	builder := flatbuffers.NewBuilder(1024)

	// Create mole weight vectors
	var moleOffsets []flatbuffers.UOffsetT
	for _, config := range configs {
		moleTypeOffset := builder.CreateString(config.MoleType)

		fbs.MoleWeightStart(builder)
		fbs.MoleWeightAddMoleType(builder, moleTypeOffset)
		fbs.MoleWeightAddWeight(builder, config.Weight)
		moleOffset := fbs.MoleWeightEnd(builder)
		moleOffsets = append(moleOffsets, moleOffset)
	}

	// Create vector
	moleVector := builder.CreateVectorOfTables(moleOffsets)

	// Create response
	fbs.GetMoleWeightRespStart(builder)
	fbs.GetMoleWeightRespAddMole(builder, moleVector)
	respOffset := fbs.GetMoleWeightRespEnd(builder)

	// Create envelope
	payloadOffset := respOffset
	fbs.EnvelopeStart(builder)
	fbs.EnvelopeAddType(builder, fbs.MessageTypeGetMoleWeightResp)
	fbs.EnvelopeAddPayload(builder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(builder)
	builder.Finish(envOffset)

	return &pb.RuntimeResponse{
		Payload: builder.FinishedBytes(),
	}, nil
}

func (s *WhackAMoleWebsocketService) GetLeaderboard(ctx context.Context, req *pb.RuntimeRequest) (*pb.RuntimeResponse, error) {
	reqEnvelope := fbs.GetRootAsEnvelope(req.Payload, 0)
	if reqEnvelope == nil {
		return nil, errors.New("invalid request envelope")
	}

	var playerID int64 = 0
	if reqEnvelope.PayloadLength() > 0 {
		reqData := fbs.GetRootAsGetLeaderboardReq(reqEnvelope.PayloadBytes(), 0)
		if reqData != nil {
			playerID = reqData.PlayerId()
		}
	}

	leaderboard, err := s.repo.GetLeaderboard(ctx, 10)
	if err != nil {
		s.log.Errorw("Failed to get leaderboard", "error", err)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	// Build flatbuffers response
	builder := flatbuffers.NewBuilder(1024)

	// Create leaderboard player vectors
	var playerOffsets []flatbuffers.UOffsetT
	var yourRank int32 = 0
	var yourScore int64 = 0

	for _, entry := range leaderboard {
		usernameOffset := builder.CreateString(entry.Username)

		fbs.LeaderboardPlayerStart(builder)
		fbs.LeaderboardPlayerAddRank(builder, entry.Rank)
		fbs.LeaderboardPlayerAddPlayerId(builder, entry.PlayerID)
		fbs.LeaderboardPlayerAddUsername(builder, usernameOffset)
		fbs.LeaderboardPlayerAddScore(builder, entry.Score)
		playerOffset := fbs.LeaderboardPlayerEnd(builder)
		playerOffsets = append(playerOffsets, playerOffset)

		// Check if this is the requested player
		if entry.PlayerID == playerID {
			yourRank = entry.Rank
			yourScore = entry.Score
		}
	}

	// Create vector
	playersVector := builder.CreateVectorOfTables(playerOffsets)

	// Create response
	fbs.GetLeaderboardRespStart(builder)
	fbs.GetLeaderboardRespAddTopPlayer(builder, playersVector)
	fbs.GetLeaderboardRespAddYourRank(builder, yourRank)
	fbs.GetLeaderboardRespAddYourScore(builder, yourScore)
	respOffset := fbs.GetLeaderboardRespEnd(builder)

	// Create envelope
	payloadOffset := respOffset
	fbs.EnvelopeStart(builder)
	fbs.EnvelopeAddType(builder, fbs.MessageTypeGetLeaderboardResp)
	fbs.EnvelopeAddPayload(builder, payloadOffset)
	envOffset := fbs.EnvelopeEnd(builder)
	builder.Finish(envOffset)

	return &pb.RuntimeResponse{
		Payload: builder.FinishedBytes(),
	}, nil
}
