package services

import (
	"fmt"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/transaction"
)

type ClickService struct {
	userChestRepository repository.UserChestRepository
	transaction transaction.TransactionManager
}

func NewClickService(
	userChestRepository repository.UserChestRepository, 
	transaction transaction.TransactionManager,
	) *ClickService {
	return &ClickService{
		userChestRepository: userChestRepository,
		transaction: transaction,
	}
}

func (s *ClickService) Damage(user *models.User,count uint) error {
	if user.UserChest.CurrentHealth <= 0 {
		return nil
	}

	s.transaction.RunInTransaction(func() error {
		err := s.userChestRepository.DecrementHealth(&user.UserChest, count)

   	if err != nil {
		return fmt.Errorf("ClickService::Damage %w", err)
	}

	if user.UserChest.CurrentHealth <= 0 {
		
	}

	return nil	
	})

	

	return nil
}
