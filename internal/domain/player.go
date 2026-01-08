package domain

type Player struct {
	ID       int64  `json:"id" gorm:"primaryKey"`
	UserName string `json:"username" gorm:"unique;not null"`
}
