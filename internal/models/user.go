package models

import (
	"time"
)

type User struct {
	ID         uint      `json:"id" gorm:"primaryKey"`               // PRIMARY KEY
	TelegramID uint64    `json:"telegram_id" gorm:"unique;not null"` // UNIQUE и NOT NULL
	Username   string    `json:"username" gorm:"size:255"`
	Session    string    `json:"session" gorm:"size:64"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"` // TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"` // TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	UserChest  UserChest `json:"user_chest" gorm:"foreignKey:UserID"`
	Balance    Balance   `json:"balance" gorm:"foreignKey:UserID"`
}
