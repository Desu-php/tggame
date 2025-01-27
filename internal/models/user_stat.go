package models

import "time"

type UserStat struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"unique;not null;constraint:OnDelete:CASCADE" json:"-"`
	Damage         uint      `gorm:"default:1" json:"damage"`
	CriticalDamage uint      `gorm:"default:0" json:"critical_damage"`
	CriticalChance float64   `gorm:"type:decimal(5,2);default:0" json:"critical_chance"`
	GoldMultiplier float64   `gorm:"type:decimal(5,2);default:0" json:"gold_multiplier"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
