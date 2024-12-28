package models

import (
	"time"
)

type Rarity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`                  // Уникальный идентификатор
	Name        string    `gorm:"type:varchar(50);not null" json:"name"` // Название сундука
	DropWeight  int       `gorm:"not null" json:"drop_weight"`           // Здоровье сундука
	Description string    `gorm:"not null" json:"description"`           // Является ли сундук по умолчанию
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`      // Время создания
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`      // Время обновления
	Color       string    `gorm:"type:varchar(255);not null" json:"color"`
	Sort        int       `gorm:"not null" json:"sort"`
}
