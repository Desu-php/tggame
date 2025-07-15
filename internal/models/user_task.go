package models

import (
	"time"
)

type UserTask struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	UserID      uint `gorm:"not null"`
	TaskID      uint `gorm:"not null"`
	Progress    uint
	CompletedAt *time.Time
	Date        time.Time
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Task        *Task
	IsNotified  bool
}

func (u *UserTask) TableName() string {
	return "user_tasks"
}

func (u *UserTask) ModelID() int {
	return int(u.ID)
}
