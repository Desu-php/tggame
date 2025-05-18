package models

import (
	"time"
)

type UserStatHistory struct {
	ID               uint      `gorm:"primaryKey"`
	UserID           uint      `gorm:"not null;index"`
	Damage           int       `gorm:"not null"`
	CriticalDamage   int       `gorm:"not null"`
	CriticalChance   float64   `gorm:"type:decimal(5,2);not null"`
	GoldMultiplier   float64   `gorm:"type:decimal(5,2);not null"`
	PassiveDamage    int       `gorm:"not null"`
	IsUpgrade        bool      `gorm:"not null"`
	AttributableType string    `gorm:"type:varchar(255);not null"`
	AttributableID   uint      `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
