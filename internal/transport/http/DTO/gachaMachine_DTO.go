package dto

type CreateGachaPlayerReq struct {
	PlayerID int64  `json:"playerID" binding:"required"`
	UserName string `json:"userName" binding:"required"`
	Coin     int64  `json:"coin" binding:"required,min=0"`
	Diamond  int64  `json:"diamond" binding:"required,min=0"`
}

type AdjustGachaPlayerCoinRequest struct {
	PlayerID int64  `json:"playerID" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=plus minus"`
}

type AdjustGachaPlayerDiamondRequest struct {
	PlayerID int64  `json:"playerID" binding:"required"`
	Amount   int64  `json:"amount" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=plus minus"`
}

type CreateGachaItemsRequest struct {
	GachaItems []CreateGachaItemRequest `json:"gachaItems" binding:"required"`
}

type CreateGachaItemRequest struct {
	Name       string `json:"name" binding:"required"`
	Rarity     string `json:"rarity" binding:"required"`
	PullWeight int32  `json:"pullWeight" binding:"required"`
}

type CreateGachaMachineRequest struct {
	Name          string                          `json:"name" binding:"required"`
	Price         int64                           `json:"price" binding:"required"`
	PriceTimesTen int64                           `json:"priceTimesTen" binding:"required"`
	SuperRarePity int32                           `json:"superRarePity" binding:"required"`
	UltraRarePity int32                           `json:"ultraRarePity" binding:"required"`
	GachaItems    []CreateGachaMachineItemRequest `json:"gachaItems"`
}

type CreateGachaMachineItemRequest struct {
	ItemID int64 `json:"itemID" binding:"required"`
}

type GetPullResultRequest struct {
	MachineID int64 `json:"machineID" binding:"required"`
	PlayerID  int64 `json:"playerID" binding:"required"`
	PullCount int32 `json:"pullCount" binding:"required,oneof=1 10"`
}
