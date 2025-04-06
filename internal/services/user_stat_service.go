package services

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/db"
	"fmt"
)

type UserStatService struct {
	db                 *db.DB
	userStatRepository repository.UserStatRepository
}

type UserStatUpgradeDto struct {
	Damage         uint
	CriticalDamage uint
	CriticalChance float64
	GoldMultiplier float64
	User           *models.User
}

func NewUserStatService(
	db *db.DB,
	userStatRepository repository.UserStatRepository,
) *UserStatService {
	return &UserStatService{
		db:                 db,
		userStatRepository: userStatRepository,
	}
}

func (s *UserStatService) Upgrade(ctx context.Context, dto UserStatUpgradeDto) error {
	userStat, err := s.userStatRepository.GetStat(ctx, dto.User)
	if err != nil {
		return err
	}

	userStat.CriticalDamage += dto.CriticalDamage
	userStat.Damage += dto.Damage
	userStat.GoldMultiplier += dto.GoldMultiplier
	userStat.CriticalChance += dto.CriticalChance

	if err := s.db.WithContext(ctx).Model(&userStat).Updates(userStat).Error; err != nil {
		return fmt.Errorf("UserStatService::Upgrade %w", err)
	}

	return nil
}
