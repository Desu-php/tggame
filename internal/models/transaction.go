package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type Transaction struct {
	ID         uint            `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     int             `gorm:"not null" json:"user_id"`
	Amount     int64           `gorm:"not null" json:"amount"`
	ModelType  string          `gorm:"type:varchar(255);not null" json:"model_type"`
	ModelID    int             `gorm:"not null" json:"model_id"`
	Type       TransactionType `gorm:"not null" json:"type"`
	OldBalance int64           `gorm:"not null" json:"old_balance"`
	NewBalance int64           `gorm:"not null" json:"new_balance"`
	CreatedAt  time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type TransactionType int16

// Возможные значения для TransactionType
const (
	TransactionTypeIncome  TransactionType = 1
	TransactionTypeExpense TransactionType = 2
)

// String возвращает строковое представление enum (value receiver)
func (t *TransactionType) String() string {
	switch *t {
	case TransactionTypeIncome:
		return "income"
	case TransactionTypeExpense:
		return "expense"
	default:
		return "unknown"
	}
}

// Scan реализует интерфейс Scanner для работы с базой данных (pointer receiver)
func (t *TransactionType) Scan(value interface{}) error {
	if value == nil {
		return errors.New("transaction type cannot be null")
	}

	intValue, ok := value.(int64)
	if !ok {
		return fmt.Errorf("expected int64, got %T", value)
	}

	*t = TransactionType(intValue)
	return nil
}

// Value реализует интерфейс Valuer для работы с базой данных (value receiver)
func (t *TransactionType) Value() (driver.Value, error) {
	return int64(*t), nil
}
