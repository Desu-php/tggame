package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
}

type transactionRepository struct {
	db *db.DB
}

func NewTransactionRepository(db *db.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	err := r.db.WithContext(ctx).Create(transaction).Error

	if err != nil {
		return fmt.Errorf("transactionRepository.Create: %w", err)
	}

	return nil
}
