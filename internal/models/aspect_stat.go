package models

import "time"

type AspectStat struct {
	ID                 uint
	AspectID           uint
	StartLevel         uint
	EndLevel           uint
	Damage             uint
	CriticalDamage     uint
	CriticalChance     float64
	GoldMultiplier     float64
	Amount             uint
	AmountGrowthFactor float64
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (a *AspectStat) TableName() string {
	return "aspect_stats"
}

func (a *AspectStat) ModelID() int {
	return int(a.ID)
}
