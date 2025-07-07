package services

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	services "example.com/v2/internal/services/item"
	"example.com/v2/pkg/transaction"
	"example.com/v2/pkg/utils"
	"fmt"
)

type UserChestService struct {
	userChestRepo        repository.UserChestRepository
	chestRepo            repository.ChestRepository
	transaction          transaction.TransactionManager
	userChestHistoryRepo repository.UserChestHistoryRepository
	itemService          *services.ItemService
	userItemService      *UserItemService
	balanceService       *BalanceService
	userStatRepo         repository.UserStatRepository
}

func NewUserChestService(
	userChestRepo repository.UserChestRepository,
	chestRepo repository.ChestRepository,
	transaction transaction.TransactionManager,
	userChestHistoryRepo repository.UserChestHistoryRepository,
	itemService *services.ItemService,
	userItemService *UserItemService,
	balanceService *BalanceService,
	userStatRepo repository.UserStatRepository,
) *UserChestService {
	return &UserChestService{
		userChestRepo:        userChestRepo,
		chestRepo:            chestRepo,
		transaction:          transaction,
		userChestHistoryRepo: userChestHistoryRepo,
		itemService:          itemService,
		userItemService:      userItemService,
		balanceService:       balanceService,
		userStatRepo:         userStatRepo,
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
		CurrentHealth: int64(chest.Health),
		Health:        chest.Health,
		Level:         1,
		Amount:        chest.Amount,
	})

	userChest.Chest = *chest

	if err != nil {
		return nil, fmt.Errorf("UserChestService::Create err %w", err)
	}

	return userChest, nil
}

func (s *UserChestService) LevelUp(ctx context.Context, userChest *models.UserChest, user *models.User) error {
	err := s.transaction.RunInTransaction(ctx, func(ctx context.Context) error {
		userChestHistory, err := s.userChestHistoryRepo.Create(ctx, userChest)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		item, err := s.itemService.GetRandomItem(ctx, &userChest.Chest.Rarity)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		_, err = s.userItemService.SetUserItem(ctx, user, item, userChestHistory)

		if err != nil {
			return fmt.Errorf("UserChestService::LevelUp %w", err)
		}

		err = s.replenish(ctx, userChest, user)

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
		uc.Health = nextChest.Health
		uc.ChestID = nextChest.ID
		uc.Chest = *nextChest
		uc.Amount = nextChest.Amount
	} else {
		s.IncreaseHealth(uc)
		s.IncreaseAmount(uc)
	}

	uc.CurrentHealth = int64(uc.Health)

	err = s.userChestRepo.Update(ctx, uc)

	if err != nil {
		return fmt.Errorf("UserChestService::Upgrade %w", err)
	}

	return nil
}

func (s *UserChestService) IncreaseHealth(uc *models.UserChest) {
	uc.Health = uint64(utils.GrowthIncrease(float64(uc.Health), uc.Chest.GrowthFactor))
}

func (s *UserChestService) IncreaseAmount(uc *models.UserChest) {
	uc.Amount = uint32(utils.GrowthIncrease(float64(uc.Amount), uc.Chest.AmountGrowthFactor))
}

func (s *UserChestService) replenish(ctx context.Context, uc *models.UserChest, user *models.User) error {
	if uc.Amount == 0 {
		return nil
	}

	userStat, err := s.userStatRepo.GetStat(ctx, user)

	if err != nil {
		return fmt.Errorf("UserChestService::replenish err %w", err)
	}

	amount := int64(uc.Amount)

	if userStat.GoldMultiplier > 0 {
		amount = int64(utils.GrowthIncrease(float64(amount), userStat.GoldMultiplier))
	}

	err = s.balanceService.Replenish(ctx, &TransactionDto{
		Amount: amount,
		User:   user,
		Model:  uc,
		Type:   models.TransactionTypeIncome,
	})

	if err != nil {
		return fmt.Errorf("UserChestService::replenish %w", err)
	}

	return nil
}
