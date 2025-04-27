package services

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/transaction"
	"fmt"
)

type UserItemService struct {
	userItemRepository repository.UserItemRepository
	UserStatService    *UserStatService
	trx                transaction.TransactionManager
}

func NewUserItemService(
	userItemRepository repository.UserItemRepository,
	userStatService *UserStatService,
	trx transaction.TransactionManager,
) *UserItemService {
	return &UserItemService{
		userItemRepository: userItemRepository,
		UserStatService:    userStatService,
		trx:                trx,
	}
}

func (s *UserItemService) SetUserItem(ctx context.Context, user *models.User, item *models.Item, userChestHistory *models.UserChestHistory) (*models.UserItem, error) {

	exists, err := s.userItemRepository.Exists(ctx, user.ID, item.ID)

	if err != nil {
		return nil, fmt.Errorf("UserItemService::SetUserItem %w", err)
	}

	var userItem *models.UserItem

	err = s.trx.RunInTransaction(ctx, func(ctx context.Context) error {
		if exists == false {
			err = s.UserStatService.Upgrade(ctx, UserStatUpgradeDto{
				Damage:         item.Damage,
				CriticalDamage: item.CriticalDamage,
				CriticalChance: item.CriticalChance,
				GoldMultiplier: item.GoldMultiplier,
				PassiveDamage:  item.PassiveDamage,
				User:           user,
			})

			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}

		userItem, err = s.userItemRepository.SetUserItem(ctx, user.ID, item, userChestHistory)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("UserItemService::SetUserItem %w", err)
	}

	return userItem, err
}
