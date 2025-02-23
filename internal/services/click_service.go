package services

import (
	"context"
	"fmt"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/transaction"
)

type ClickService struct {
	userChestRepository repository.UserChestRepository
	transaction         transaction.TransactionManager
	userChestService    *UserChestService
}

func NewClickService(
	userChestRepository repository.UserChestRepository,
	transaction transaction.TransactionManager,
	userChestService *UserChestService,
) *ClickService {
	return &ClickService{
		userChestRepository: userChestRepository,
		transaction:         transaction,
		userChestService:    userChestService,
	}
}

func (s *ClickService) Damage(ctx context.Context, user *models.User, count uint) error {

	userChest, err := s.userChestRepository.FindByUser(ctx, user)

	if err != nil {
		return fmt.Errorf("ClickService::Damage %w", err)
	}

	user.UserChest = *userChest

	err = s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		err := s.userChestRepository.DecrementHealth(ctx, &user.UserChest, count)

		if err != nil {
			return fmt.Errorf("ClickService::Damage %w", err)
		}

		if user.UserChest.CurrentHealth <= 0 {
			err = s.userChestService.LevelUp(ctx, &user.UserChest, user)

			if err != nil {
				return fmt.Errorf("ClickService::Damage %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
