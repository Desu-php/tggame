package services

import (
	"context"
	"fmt"
	"math"

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
	userItemService      *UserItemService
}

func NewUserChestService(
	userChestRepo repository.UserChestRepository,
	chestRepo repository.ChestRepository,
	transaction transaction.TransactionManager,
	userChestHistoryRepo repository.UserChestHistoryRepository,
	itemService *services.ItemService,
	userItemService *UserItemService,
) *UserChestService {
	return &UserChestService{
		userChestRepo:        userChestRepo,
		chestRepo:            chestRepo,
		transaction:          transaction,
		userChestHistoryRepo: userChestHistoryRepo,
		itemService:          itemService,
		userItemService:      userItemService,
	}
}

func (s *UserChestService) Create(ctx context.Context, user *models.User) (*models.UserChest, error) {
	chest, err := s.chestRepo.GetDefault(ctx)

	if err != nil {
		return nil, fmt.Errorf("UserChestService::Create err %w", err)
	}

	if chest == nil {
		return nil, fmt.Errorf("UserChestService::Create default chest not found")
	}

	userChest, err := s.userChestRepo.Create(ctx, &models.UserChest{
		UserID:        user.ID,
		ChestID:       chest.ID,
		CurrentHealth: chest.Health,
		Health:        uint(chest.Health),
		Level:         1,
	})

	userChest.Chest = *chest

	if err != nil {
		return nil, fmt.Errorf("UserChestService::Create err %w", err)
	}

	return userChest, nil
}

func (s *UserChestService) LevelUp(ctx context.Context, userChest *models.UserChest) error {
	err := s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		userChestHistory, err := s.userChestHistoryRepo.Create(ctx, userChest)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		item, err := s.itemService.GetRandomItem(ctx)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		_, err = s.userItemService.SetUserItem(ctx, userChest.UserID, item, userChestHistory)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		err = s.Upgrade(ctx, userChest)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *UserChestService) Upgrade(ctx context.Context, uc *models.UserChest) error {
	uc.Level++

	nextChest, err := s.chestRepo.GetNextChest(ctx, uint(uc.Level))

	if err != nil {
		return fmt.Errorf("UserChestService::Upgrade err %w", err)
	}

	if nextChest.ID != uc.ChestID {
		uc.Health = uint(nextChest.Health)
		uc.ChestID = nextChest.ID
		uc.Chest = *nextChest
	} else {
		s.IncreaseHealth(uc)
	}

	uc.CurrentHealth = int(uc.Health)

	err = s.userChestRepo.Update(ctx, uc)

	if err != nil {
		return fmt.Errorf("UserChestService::Upgrade %w", err)
	}

	return nil
}

func (s *UserChestService) IncreaseHealth(uc *models.UserChest) {
	increase := float64(uc.Health) * (uc.Chest.GrowthFactor / 100)
	uc.Health = uint(math.Round(float64(uc.Health) + increase))
}
