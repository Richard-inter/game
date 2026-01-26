package domain

type WhackAMolePlayer struct {
	Player Player `gorm:"embedded;embeddedPrefix:player_"`
}

type MoleWeightConfig struct {
	ID       int64  `gorm:"column:id;primaryKey" json:"id"`
	MoleType string `gorm:"column:mole_type" json:"moleType"`
	Weight   int32  `gorm:"column:weight" json:"weight"`
}

type LeaderBoard struct {
	PlayerID int64  `gorm:"column:player_id;primaryKey;index:idx_leaderboard_score,priority:2" json:"playerID"`
	Username string `gorm:"column:username" json:"username"`
	Score    int64  `gorm:"column:score;index:idx_leaderboard_score,priority:1,sort:desc" json:"score"`
	Rank     int32  `gorm:"column:rank;index:idx_leaderboard_rank" json:"rank"`
}

func (LeaderBoard) TableName() string {
	return "whackAMole_leaderboard"
}

func (WhackAMolePlayer) TableName() string {
	return "whackAMole_player"
}

func (MoleWeightConfig) TableName() string {
	return "whackAMole_mole_weight_config"
}
