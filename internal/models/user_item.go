package models

import (
	"gorm.io/gorm"
	"time"
)

type UserItem struct {
	ID                 uint             `gorm:"primaryKey" json:"id"`
	UserID             uint             `gorm:"not null" json:"user_id"`
	ItemID             uint             `gorm:"not null" json:"item_id"`
	UserChestHistoryID uint             `gorm:"not null" json:"user_chest_history_id"`
	CreatedAt          time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	Item               Item             `gorm:"foreignKey:ItemID" json:"item"`
	User               User             `gorm:"foreignKey:UserID" json:"user"`
	UserChestHistory   UserChestHistory `gorm:"foreignKey:UserChestHistoryID" json:"user_chest_history"`
	DeletedAt          gorm.DeletedAt
}
