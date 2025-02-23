package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
)

type BalanceRepository interface {
	Create(ctx context.Context, user *models.User) (*models.Balance, error)
	Update(ctx context.Context, balance *models.Balance) (*models.Balance, error)
	FindBalance(ctx context.Context, user *models.User) (*models.Balance, error)
}

type BalanceRepositoryImpl struct {
	db *db.DB
}

func NewBalanceRepository(db *db.DB) BalanceRepository {
	return &BalanceRepositoryImpl{db: db}
}

func (b *BalanceRepositoryImpl) Create(ctx context.Context, user *models.User) (*models.Balance, error) {
	balance := &models.Balance{
		UserID:  user.ID,
		Balance: 0,
	}

	err := b.db.WithContext(ctx).Create(balance).Error

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::Create %w", err)
	}

	return balance, nil
}

func (b *BalanceRepositoryImpl) FindBalance(ctx context.Context, user *models.User) (*models.Balance, error) {
	balance := &models.Balance{
		UserID: user.ID,
	}

	err := b.db.WithContext(ctx).First(balance).Error

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::SetBalance %w", err)
	}

	return balance, nil
}

func (b *BalanceRepositoryImpl) Update(ctx context.Context, balance *models.Balance) (*models.Balance, error) {
	err := b.db.WithContext(ctx).Save(balance).Error

	if err != nil {
		return nil, fmt.Errorf("BalanceRepository::Update %w", err)
	}

	return balance, nil
}
