package dto

// CreatePlayerRequest represents the HTTP request for creating a player
type CreatePlayerRequest struct {
	UserName string `json:"userName" binding:"required"`
}
