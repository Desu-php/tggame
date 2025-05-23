package models

import (
	"time"
)

type UserChestHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`                     // Уникальный идентификатор
	UserChestID uint      `gorm:"not null" json:"user_id"`                  // Ссылка на пользователя
	Health      uint64    `gorm:"not null" json:"health"`                   // Текущее здоровье сундука
	Level       int       `gorm:"default:1" json:"level"`                   // Уровень сундука
	Amount      uint32    `json:"amount"`                                   // Уровень сундука
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`         // Время создания
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`         // Время обновления
	UserChest   UserChest `gorm:"foreignKey:UserChestID" json:"user_chest"` // Связь с сундуком
}
