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
	PassiveDamage  uint
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
	userStat.PassiveDamage += dto.PassiveDamage

	err = s.db.WithContext(ctx).Model(userStat).Updates(userStat).Error

	if err != nil {
		return fmt.Errorf("UserStatService::Upgrade %w", err)
	}

	return nil
}

func (s *UserStatService) Downgrade(ctx context.Context, dto UserStatUpgradeDto) error {
	userStat, err := s.userStatRepository.GetStat(ctx, dto.User)
	if err != nil {
		return err
	}

	userStat.CriticalDamage -= dto.CriticalDamage
	userStat.Damage -= dto.Damage
	userStat.GoldMultiplier -= dto.GoldMultiplier
	userStat.CriticalChance -= dto.CriticalChance
	userStat.PassiveDamage -= dto.PassiveDamage

	err = s.db.WithContext(ctx).Model(userStat).Updates(map[string]interface{}{
		"critical_damage": userStat.CriticalDamage,
		"damage":          userStat.Damage,
		"gold_multiplier": userStat.GoldMultiplier,
		"critical_chance": userStat.CriticalChance,
		"passive_damage":  userStat.PassiveDamage,
	}).Error

	if err != nil {
		return fmt.Errorf("UserStatService::Downgrade %w", err)
	}

	return nil
}
