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
	Damage         int        `gorm:"not null"`
	CriticalDamage int        `gorm:"not null"`
	CriticalChance float64    `gorm:"type:decimal(5,2);not null"`
	GoldMultiplier float64    `gorm:"type:decimal(5,2);not null"`
	Amount         int        `gorm:"not null"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}
