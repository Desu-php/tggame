package transaction

import (
	"context"
	"example.com/v2/pkg/db"
	"gorm.io/gorm"
)

type GormTransactionManager struct {
	db *db.DB
}

func NewTransactionManager(db *db.DB) TransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) RunInTransaction(ctx context.Context, fn TransactionFunc) error {
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tm.db.PutTxToContext(ctx, tx))
	})
}
