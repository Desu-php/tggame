package models

import (
	"errors"
	"fmt"
	"time"
)

type Aspect struct {
	ID          uint
	Name        string
	Type        AspectType
	Description string
	Image       string
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type AspectType string

// Возможные значения для TransactionType
const (
	Aspects AspectType = "aspect"
	Booster AspectType = "booster"
)

func (a *AspectType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("aspect type cannot be null")
	}

	stringValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string, got %T", value)
	}

	*a = AspectType(stringValue)
	return nil
}

func (a *Aspect) IsAspect() bool {
	return a.Type == Aspects
}

func (a *Aspect) IsBooster() bool {
	return a.Type == Booster
}
