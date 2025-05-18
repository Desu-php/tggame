package models

import (
	"time"
)

type UserTask struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	UserID      uint `gorm:"not null"`
	TaskID      uint `gorm:"not null"`
	CompletedAt *time.Time
}
