package models

import (
	"time"
)

type ItemType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`                  // Уникальный идентификатор
	Name        string    `gorm:"type:varchar(50);not null" json:"name"` // Название сундука
	Description string    `gorm:"not null" json:"description"`           // Является ли сундук по умолчанию
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`      // Время создания
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // Время обновления
}
