package models

import (
	"time"
)

type Chest struct {
	ID           uint      `gorm:"primaryKey" json:"id"`                            // Уникальный идентификатор
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`          // Название сундука
	Health       int       `gorm:"not null" json:"health"`                          // Здоровье сундука
	IsDefault    bool      `gorm:"not null" json:"is_default"`                      // Является ли сундук по умолчанию
	GrowthFactor float64   `gorm:"type:numeric(5,2);not null" json:"growth_factor"` // Коэффициент роста
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`                // Время создания
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`                // Время обновления
}
