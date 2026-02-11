package domain

import (
	"time"

	"gorm.io/gorm"
)

type GachaPlayer struct {
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

type GachaMachine struct {
	ID            int64          `gorm:"column:id;primaryKey" json:"machineID"`
	Name          string         `gorm:"column:name" json:"name"`
	Price         int64          `gorm:"column:price" json:"price"`
	PriceTimesTen int64          `gorm:"column:price_times_ten" json:"priceTimesTen"`
	SuperRarePity int32          `gorm:"column:super_rare_pity" json:"superRarePity"`
	UltraRarePity int32          `gorm:"column:ultra_rare_pity" json:"ultraRarePity"`
	IsActive      bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy     string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy     string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy     *string        `gorm:"column:deleted_by" json:"deletedBy"`

	Items []GachaMachineItem `gorm:"foreignKey:GachaMachineID;constraint:OnDelete:CASCADE"`
}

type GachaMachineItem struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GachaMachineID int64          `gorm:"column:gacha_machine_id" json:"gachaMachineID"`
	ItemID         int64          `gorm:"column:item_id" json:"itemID"`
	IsActive       bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy      string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy      string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy      *string        `gorm:"column:deleted_by" json:"deletedBy"`

	Item GachaItem `gorm:"foreignKey:ItemID;references:ID"`
}

type GachaItem struct {
	ID         int64          `gorm:"column:id;primaryKey" json:"itemID"`
	Name       string         `gorm:"column:name" json:"name"`
	Rarity     string         `gorm:"column:rarity" json:"rarity"`
	PullWeight int32          `gorm:"column:pull_weight" json:"pullWeight"`
	IsActive   bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy  string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy  string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy  *string        `gorm:"column:deleted_by" json:"deletedBy"`
}

type GachaPullSession struct {
	ID             int64          `gorm:"column:id;primaryKey" json:"sessionID"`
	GachaMachineID int64          `gorm:"column:gacha_machine_id" json:"gachaMachineID"`
	PlayerID       int64          `gorm:"column:player_id" json:"playerID"`
	PullCount      int32          `gorm:"column:pull_count" json:"pullCount"`
	IsActive       bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy      string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy      string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy      *string        `gorm:"column:deleted_by" json:"deletedBy"`

	GachaPullHistories []GachaPullHistory `gorm:"foreignKey:GachaPullSessionID;references:ID"`
}

type GachaPullHistory struct {
	ID                 int64          `gorm:"column:id;primaryKey" json:"historyID"`
	GachaPullSessionID int64          `gorm:"column:gacha_pull_session_id" json:"gachaPullSessionID"`
	ItemID             int64          `gorm:"column:item_id" json:"itemID"`
	IsActive           bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt          time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy          string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy          string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy          *string        `gorm:"column:deleted_by" json:"deletedBy"`

	Item GachaItem `gorm:"foreignKey:ItemID;references:ID"`
}

type GachaPityState struct {
	ID                 int64 `gorm:"column:id;primaryKey" json:"id"`
	GachaMachineID     int64 `gorm:"column:gacha_machine_id" json:"gachaMachineID"`
	PlayerID           int64 `gorm:"column:player_id" json:"playerID"`
	SuperRarePityCount int32 `gorm:"column:super_rare_pity_count" json:"superRarePityCount"`
	UltraRarePityCount int32 `gorm:"column:ultra_rare_pity_count" json:"ultraRarePityCount"`
}

type GachaMachineRTPState struct {
	GachaMachineID int64          `gorm:"column:gacha_machine_id;primaryKey" json:"gachaMachineID"`
	TargetRTP      float64        `gorm:"column:target_rtp" json:"targetRTP"`
	TotalRevenue   int64          `gorm:"column:total_revenue" json:"totalRevenue"`
	TotalPlays     int64          `gorm:"column:total_plays" json:"totalPlays"`
	TotalPayout    int64          `gorm:"column:total_payout" json:"totalPayout"`
	IsActive       bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy      string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy      string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy      *string        `gorm:"column:deleted_by" json:"deletedBy"`
}

type GachaPlayerInventory struct {
	PlayerID  int64          `gorm:"column:player_id;primaryKey" json:"playerID"`
	ItemID    int64          `gorm:"column:item_id;primaryKey" json:"itemID"`
	Quantity  int32          `gorm:"column:quantity" json:"quantity"`
	IsActive  bool           `gorm:"column:is_active;default:true" json:"isActive"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	CreatedBy string         `gorm:"column:created_by" json:"createdBy"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	UpdatedBy string         `gorm:"column:updated_by" json:"updatedBy"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt"`
	DeletedBy *string        `gorm:"column:deleted_by" json:"deletedBy"`
}

func (GachaPlayer) TableName() string {
	return "gacha_player"
}

func (GachaMachine) TableName() string {
	return "gacha_machine"
}

func (GachaMachineItem) TableName() string {
	return "gacha_machine_item"
}

func (GachaItem) TableName() string {
	return "gacha_item"
}

func (GachaPullSession) TableName() string {
	return "gacha_pull_session"
}

func (GachaPullHistory) TableName() string {
	return "gacha_pull_history"
}

func (GachaPityState) TableName() string {
	return "gacha_pity_state"
}

func (GachaMachineRTPState) TableName() string {
	return "gacha_machine_rtp_state"
}

func (GachaPlayerInventory) TableName() string {
	return "gacha_player_inventory"
}
