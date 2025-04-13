package models

import "time"

type UserAspect struct {
	ID             uint       `gorm:"primaryKey"`
	UserID         uint       `gorm:"not null;index"`
	User           User       `gorm:"foreignKey:UserID"`
	AspectID       uint       `gorm:"not null;index"`
	Aspect         Aspect     `gorm:"foreignKey:AspectID"`
	AspectStatID   uint       `gorm:"not null;index"`
	AspectStat     AspectStat `gorm:"foreignKey:AspectStatID"`
	Level          int        `gorm:"not null"`
	Damage         uint       `gorm:"not null"`
	CriticalDamage uint       `gorm:"not null"`
	CriticalChance float64    `gorm:"type:decimal(5,2);not null"`
	GoldMultiplier float64    `gorm:"type:decimal(5,2);not null"`
	Amount         uint       `gorm:"not null"`
	PassiveDamage  uint       `gorm:"default:0" json:"passive_damage"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}
