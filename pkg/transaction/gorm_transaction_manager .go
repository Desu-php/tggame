package transaction

import "gorm.io/gorm"

type GormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &GormTransactionManager{db: db}
}

func (tm *GormTransactionManager) RunInTransaction(fn TransactionFunc) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		return fn()
	})
}