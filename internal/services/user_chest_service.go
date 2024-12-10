package services

import (
	"fmt"

	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	services "example.com/v2/internal/services/item"
	"example.com/v2/pkg/transaction"
)

type UserChestService struct {
	userChestRepo        repository.UserChestRepository
	chestRepo            repository.ChestRepository
	transaction          transaction.TransactionManager
	userChestHistoryRepo repository.UserChestHistoryRepository
	itemService          *services.ItemService
}

func NewUserChestService(
	userChestRepo repository.UserChestRepository,
	chestRepo repository.ChestRepository,
	transaction transaction.TransactionManager,
	userChestHistoryRepo repository.UserChestHistoryRepository,
	itemService *services.ItemService,
) *UserChestService {
	return &UserChestService{
		userChestRepo:        userChestRepo,
		chestRepo:            chestRepo,
		transaction:          transaction,
		userChestHistoryRepo: userChestHistoryRepo,
		itemService:          itemService,
	}
}

func (s *UserChestService) Create(user *models.User) (*models.UserChest, error) {
	chest, err := s.chestRepo.GetDefault()

	if err != nil {
		return nil, fmt.Errorf("UserChestService::Create err %w", err)
	}

	if chest == nil {
		return nil, fmt.Errorf("UserChestService::Create default chest not found")
	}

	userChest, err := s.userChestRepo.Create(&models.UserChest{
		UserID:        user.ID,
		ChestID:       chest.ID,
		CurrentHealth: chest.Health,
		Level:         1,
	})

	if err != nil {
		return nil, fmt.Errorf("UserChestService::Create err %w", err)
	}

	userChest.Chest = *chest

	return userChest, nil
}

func (s *UserChestService) LevelUp(userChest *models.UserChest) (*models.UserChest, error) {
	s.transaction.RunInTransaction(func() error {
		_, err := s.userChestHistoryRepo.Create(userChest)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		item, err := s.itemService.GetRandomItem()

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

	})
}
