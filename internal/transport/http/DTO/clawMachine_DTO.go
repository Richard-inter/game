package dto

// CreateClawMachineRequest represents the HTTP request for creating a claw machine
type CreateClawMachineRequest struct {
	Name    string                         `json:"name" binding:"required"`
	Price   int64                          `json:"price" binding:"required"`
	MaxItem int32                          `json:"maxItem" binding:"required"`
	Items   []CreateClawMachineItemRequest `json:"items"`
}

// CreateClawMachineItemRequest represents an item in the claw machine creation request
type CreateClawMachineItemRequest struct {
	ItemID int64 `json:"itemID" binding:"required"`
}

// CreateClawItemsRequest represents the HTTP request for creating claw items
type CreateClawItemsRequest struct {
	ClawItems []CreateClawItemRequest `json:"clawItems" binding:"required"`
}

// CreateClawItemRequest represents the HTTP request for creating a single claw item
type CreateClawItemRequest struct {
	Name            string `json:"name" binding:"required"`
	Rarity          string `json:"rarity" binding:"required"`
	SpawnPercentage int64  `json:"spawnPercentage" binding:"required"`
	CatchPercentage int64  `json:"catchPercentage" binding:"required"`
	MaxItemSpawned  int64  `json:"maxItemSpawned" binding:"required"`
}

type CreateClawPlayerRequest struct {
	PlayerID int64  `json:"playerID" binding:"required"`
	UserName string `json:"userName" binding:"required"`
	Coin     int64  `json:"coin" binding:"required"`
	Diamond  int64  `json:"diamond" binding:"required"`
}
