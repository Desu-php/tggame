package models

import (
	"time"
)

type Item struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Image       string    `gorm:"type:varchar(250);not null" json:"image"`
	TypeID      uint      `gorm:"not null" json:"type_id"`
	RarityID    uint      `gorm:"not null" json:"rarity_id"`
	Description string    `gorm:"not null" json:"description"`
	DropChance  float32   `gorm:"not null" json:"drop_change"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Type        ItemType  `gorm:"foreignKey:TypeID" json:"type"`
	Rarity      Rarity    `gorm:"foreignKey:RarityID" json:"rarity"`
}
