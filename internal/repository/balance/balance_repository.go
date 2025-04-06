package repository

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"fmt"
	"gorm.io/gorm/clause"
)

type BalanceRepository interface {
	Create(ctx context.Context, user *models.User) (*models.Balance, error)
	Update(ctx context.Context, balance *models.Balance) (*models.Balance, error)
	FindBalance(ctx context.Context, user *models.User) (*models.Balance, error)
	LockForUpdate() *BalanceRepositoryImpl
}

type BalanceRepositoryImpl struct {
	db            *db.DB
	lockForUpdate bool
}

func NewBalanceRepository(db *db.DB) BalanceRepository {
	return &BalanceRepositoryImpl{db: db, lockForUpdate: false}
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

func (b *BalanceRepositoryImpl) LockForUpdate() *BalanceRepositoryImpl {
	b.lockForUpdate = true

	return b
}

func (b *BalanceRepositoryImpl) FindBalance(ctx context.Context, user *models.User) (*models.Balance, error) {
	balance := &models.Balance{}

	query := b.db.WithContext(ctx).Where("user_id = ?", user.ID)

	if b.lockForUpdate {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
		b.lockForUpdate = false
	}

	err := query.First(balance).Error

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
