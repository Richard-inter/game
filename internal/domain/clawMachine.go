package domain

import (
	"time"

	"gorm.io/gorm"
)

type ClawMachine struct {
	ID        int64          `gorm:"column:id;primaryKey" json:"machineID"`
	Name      string         `gorm:"column:name" json:"name"`
	Price     int64          `gorm:"column:price" json:"price"`
	MaxItem   int32          `gorm:"column:max_item" json:"maxItem"`
	IsActive  bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy *string        `gorm:"column:deleted_by" json:"deletedBy"`

	Items []ClawMachineItem `gorm:"foreignKey:ClawMachineID;constraint:OnDelete:CASCADE"`
}

type ClawMachineItem struct {
	ID            int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	ClawMachineID int64          `gorm:"column:claw_machine_id" json:"clawMachineID"`
	ItemID        int64          `gorm:"column:item_id" json:"itemID"`
	IsActive      bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy     string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy     string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy     *string        `gorm:"column:deleted_by" json:"deletedBy"`

	Item ClawItem `gorm:"foreignKey:ItemID;references:ID"`
}

type ClawItem struct {
	ID              int64          `gorm:"column:id;primaryKey" json:"itemID"`
	Name            string         `gorm:"column:name" json:"name"`
	Rarity          string         `gorm:"column:rarity" json:"rarity"`
	SpawnPercentage float64        `gorm:"column:spawn_percentage" json:"spawnPercentage"`
	CatchPercentage float64        `gorm:"column:catch_percentage" json:"catchPercentage"`
	MaxItemSpawned  int64          `gorm:"column:max_item_spawned" json:"maxItemSpawned"`
	IsActive        bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy       string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy       string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy       *string        `gorm:"column:deleted_by" json:"deletedBy"`
}

type ClawPlayer struct {
	Player    Player         `gorm:"embedded;embeddedPrefix:player_"`
	Coin      int64          `gorm:"column:coin;not null" json:"coin"`
	Diamond   int64          `gorm:"column:diamond;not null" json:"diamond"`
	IsActive  bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy *string        `gorm:"column:deleted_by" json:"deletedBy"`
}

type ClawMachineGameRecord struct {
	ID            int64          `gorm:"column:id;primaryKey" json:"gameID"`
	ClawMachineID int64          `gorm:"column:claw_machine_id" json:"clawMachineID"`
	PlayerID      int64          `gorm:"column:player_id" json:"playerID"`
	TouchedItemID *int64         `gorm:"column:touched_item_id" json:"touchedItemID"`
	Catched       bool           `gorm:"column:catched" json:"catched"`
	IsActive      bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy     string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy     string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy     *string        `gorm:"column:deleted_by" json:"deletedBy"`

	Machine     ClawMachine `gorm:"foreignKey:ClawMachineID;references:ID"`
	TouchedItem ClawItem    `gorm:"foreignKey:TouchedItemID;references:ID"`
}

type ClawMachineRTPState struct {
	ClawMachineID int64          `gorm:"column:claw_machine_id;primaryKey"`
	TargetRTP     float64        `gorm:"column:target_rtp"`
	TotalPlays    int64          `gorm:"column:total_plays"`
	TotalRevenue  int64          `gorm:"column:total_revenue"`
	TotalPayout   int64          `gorm:"column:total_payout"`
	IsActive      bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy     string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy     string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy     *string        `gorm:"column:deleted_by" json:"deletedBy"`
}

func (ClawMachine) TableName() string {
	return "claw_machine"
}

func (ClawMachineItem) TableName() string {
	return "claw_machine_item"
}

func (ClawItem) TableName() string {
	return "claw_item"
}

func (ClawPlayer) TableName() string {
	return "claw_player"
}

func (ClawMachineGameRecord) TableName() string {
	return "claw_machine_game_record"
}

func (ClawMachineRTPState) TableName() string {
	return "claw_machine_rtp_state"
}
