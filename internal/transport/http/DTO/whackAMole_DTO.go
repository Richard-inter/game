package dto

type CreateWhackAMolePlayerRequest struct {
	PlayerID int64  `json:"playerId" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type GetPlayerInfoRequest struct {
	PlayerID int64 `json:"playerId" binding:"required"`
}

type GetLeaderboardRequest struct {
	Limit int32 `json:"limit,omitempty"`
}

type CreateMoleWeightConfigRequest struct {
	MoleType string `json:"moleType" binding:"required"`
	Weight   int32  `json:"weight" binding:"required,gt=0"`
}

type UpdateMoleWeightConfigRequest struct {
	ID       int64  `json:"id" binding:"required"`
	MoleType string `json:"moleType"`
	Weight   int32  `json:"weight" validate:"gt=0"`
}

type UpdateScoreRequest struct {
	PlayerID int64 `json:"playerID" binding:"required"`
	Score    int64 `json:"score" binding:"required,gt=0"`
}
