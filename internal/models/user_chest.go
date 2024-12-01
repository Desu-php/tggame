package models

import (
	"time"
)

type UserChest struct {
	ID            uint      `gorm:"primaryKey" json:"id"`             // Уникальный идентификатор
	UserID        uint      `gorm:"not null" json:"user_id"`          // Ссылка на пользователя
	ChestID       uint      `gorm:"not null" json:"chest_id"`         // Ссылка на сундук
	CurrentHealth int       `gorm:"not null" json:"current_health"`   // Текущее здоровье сундука
	Level         int       `gorm:"default:1" json:"level"`           // Уровень сундука
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"` // Время создания
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Время обновления
	Chest         Chest     `gorm:"foreignKey:ChestID" json:"chest"`  // Связь с сундуком
}
