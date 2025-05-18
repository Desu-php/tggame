package models

import (
	"errors"
	"fmt"
	"time"
)

type TaskType string

// Возможные значения для TransactionType
const (
	TaskTypeDamage    TaskType = "damage"
	TaskTypeDestroy   TaskType = "destroy"
	TaskTypeClickLink TaskType = "click_link"
)

func (t *TaskType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("transaction type cannot be null")
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	*t = TaskType(strValue)
	return nil
}

type Task struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text"`
	Type        TaskType  `gorm:"type:varchar(50);not null"`
	TargetValue int       `gorm:"not null"`
	Amount      int64     `gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
