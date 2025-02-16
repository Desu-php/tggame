package transaction

import (
	"context"
	"gorm.io/gorm"
)

type txKey struct{}

type GormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) RunInTransaction(ctx context.Context, fn TransactionFunc) error {
	// Проверяем, есть ли уже открытая транзакция в текущем контексте
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		// Если транзакция уже есть, передаем её в функцию
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	}

	// Если транзакции нет, создаем новую и передаем в контекст
	return tm.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}
