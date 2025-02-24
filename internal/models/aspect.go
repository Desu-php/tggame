package models

import "time"

type Aspect struct {
	ID          uint
	Name        string
	Description string
	Image       string
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
