package services

import (
	"context"
	"example.com/v2/internal/models"
	balance "example.com/v2/internal/repository/balance"
	"example.com/v2/pkg/transaction"
	"fmt"
	"log"
)

type BalanceService struct {
	trxRepository     balance.TransactionRepository
	balanceRepository balance.BalanceRepository
	trx               transaction.TransactionManager
}
type TransactionDto struct {
	Amount int64
	User   *models.User
	Model  Model
	Type   models.TransactionType
}

type Model interface {
	TableName() string
	ModelID() int
}

func NewBalanceService(
	trxRepository balance.TransactionRepository,
	balanceRepository balance.BalanceRepository,
	trx transaction.TransactionManager,
) *BalanceService {
	return &BalanceService{trxRepository, balanceRepository, trx}
}

func (s *BalanceService) Replenish(ctx context.Context, dto *TransactionDto) error {
	log.Println(dto)

	if dto.Amount <= 0 {
		return fmt.Errorf("BalanceService::Replenish amount must be greater than zero")
	}

	_, err := s.CreateTransaction(ctx, dto)

	if err != nil {
		return fmt.Errorf("BalanceService::Replenish amount must be greater than zero")
	}

	return nil
}

func (s *BalanceService) CreateTransaction(ctx context.Context, dto *TransactionDto) (*models.Transaction, error) {
	var createdTransaction *models.Transaction

	err := s.trx.RunInTransaction(ctx,
		func(ctx context.Context) error {
			findBalance, err := s.balanceRepository.FindBalance(ctx, dto.User)

			if err != nil {
				return fmt.Errorf("balanceService.CreateTransaction: %w", err)
			}

			oldBalance := findBalance.Balance

			findBalance.Balance = findBalance.Balance + dto.Amount

			_, err = s.balanceRepository.Update(ctx, findBalance)

			if err != nil {
				return fmt.Errorf("balanceService.CreateTransaction: %w", err)

			}

			createdTransaction = &models.Transaction{
				UserID:     dto.User.ID,
				Amount:     dto.Amount,
				ModelType:  dto.Model.TableName(),
				ModelID:    dto.Model.ModelID(),
				Type:       dto.Type,
				OldBalance: oldBalance,
				NewBalance: findBalance.Balance,
			}

			err = s.trxRepository.Create(ctx, createdTransaction)

			if err != nil {
				return fmt.Errorf("balanceService.CreateTransaction: %w", err)
			}

			return nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("balanceService.CreateTransaction: %w", err)
	}

	return createdTransaction, nil
}
