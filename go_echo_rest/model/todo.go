package model

type ToDo struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	Status      bool   `json:"status"`
	UserID      uint   `json:"user_id" gorm:"not null"`
	User        User   `json:"-" gorm:"foreignKey:UserID"`
}
