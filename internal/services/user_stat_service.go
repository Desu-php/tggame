package services

import (
	"context"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/transaction"
	"fmt"
	"time"
)

type UserStatService struct {
	db                 *db.DB
	userStatRepository repository.UserStatRepository
	trx                transaction.TransactionManager
}

type UserStatUpgradeDto struct {
	Damage         uint
	CriticalDamage uint
	CriticalChance float64
	GoldMultiplier float64
	PassiveDamage  uint
	User           *models.User
	Attributable   Attributable
}

type UserStatDowngradeDto struct {
	User         *models.User
	Attributable Attributable
}

type Attributable interface {
	AttributableID() uint
	AttributableName() string
}

func NewUserStatService(
	db *db.DB,
	userStatRepository repository.UserStatRepository,
	trx transaction.TransactionManager,
) *UserStatService {
	return &UserStatService{
		db:                 db,
		userStatRepository: userStatRepository,
		trx:                trx,
	}
}

func (s *UserStatService) Upgrade(ctx context.Context, dto *UserStatUpgradeDto) error {
	userStat, err := s.userStatRepository.GetStat(ctx, dto.User)
	if err != nil {
		return err
	}

	userStat.CriticalDamage += dto.CriticalDamage
	userStat.Damage += dto.Damage
	userStat.GoldMultiplier += dto.GoldMultiplier
	userStat.CriticalChance += dto.CriticalChance
	userStat.PassiveDamage += dto.PassiveDamage

	err = s.trx.RunInTransaction(ctx, func(ctx context.Context) error {
		err = s.db.WithContext(ctx).Model(userStat).Updates(userStat).Error

		if err != nil {
			return err
		}

		err = s.log(ctx, dto, true)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("UserStatService::Upgrade %w", err)
	}

	return nil
}

func (s *UserStatService) Downgrade(ctx context.Context, dto *UserStatDowngradeDto) error {
	userStat, err := s.userStatRepository.GetStat(ctx, dto.User)
	if err != nil {
		return err
	}

	var userStatHistory models.UserStatHistory

	err = s.db.WithContext(ctx).Model(&models.UserStatHistory{}).
		Where("user_id = ?", dto.User.ID).
		Where("attributable_type = ?", dto.Attributable.AttributableName()).
		Where("attributable_id = ?", dto.Attributable.AttributableID()).
		Where("is_upgrade = true").
		Order("id desc").
		First(&userStatHistory).Error

	if err != nil {
		return fmt.Errorf("UserStatService::Downgrade %w", err)
	}

	userStat.CriticalDamage -= uint(userStatHistory.CriticalDamage)
	userStat.Damage -= uint(userStatHistory.Damage)
	userStat.GoldMultiplier -= userStatHistory.GoldMultiplier
	userStat.CriticalChance -= userStatHistory.CriticalChance
	userStat.PassiveDamage -= uint(userStatHistory.PassiveDamage)

	err = s.trx.RunInTransaction(ctx, func(ctx context.Context) error {
		err = s.db.WithContext(ctx).Model(userStat).Updates(map[string]interface{}{
			"critical_damage": userStat.CriticalDamage,
			"damage":          userStat.Damage,
			"gold_multiplier": userStat.GoldMultiplier,
			"critical_chance": userStat.CriticalChance,
			"passive_damage":  userStat.PassiveDamage,
		}).Error

		if err != nil {
			return fmt.Errorf("UserStat updates %w", err)
		}

		err = s.log(ctx, &UserStatUpgradeDto{
			Damage:         uint(userStatHistory.Damage),
			CriticalDamage: uint(userStatHistory.CriticalDamage),
			CriticalChance: userStatHistory.CriticalChance,
			GoldMultiplier: userStatHistory.GoldMultiplier,
			PassiveDamage:  uint(userStatHistory.PassiveDamage),
			User:           dto.User,
			Attributable:   dto.Attributable,
		}, false)

		if err != nil {
			return fmt.Errorf("UserStat logs %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("UserStatService::Downgrade %w", err)
	}

	return nil
}

func (s *UserStatService) log(ctx context.Context, dto *UserStatUpgradeDto, isUpgrade bool) error {
	damage := int(dto.Damage)
	criticalDamage := int(dto.CriticalDamage)
	criticalChance := dto.CriticalChance
	goldMultiplier := dto.GoldMultiplier
	passiveDamage := int(dto.PassiveDamage)

	if isUpgrade == false {
		damage = -damage
		criticalDamage = -criticalDamage
		criticalChance = -criticalChance
		goldMultiplier = -goldMultiplier
		passiveDamage = -passiveDamage
	}

	if dto.Attributable == nil {
		return fmt.Errorf("UserStatService::log: Attributable is nil")
	}

	err := s.db.WithContext(ctx).Create(&models.UserStatHistory{
		UserID:           dto.User.ID,
		Damage:           damage,
		CriticalDamage:   criticalDamage,
		CriticalChance:   criticalChance,
		GoldMultiplier:   goldMultiplier,
		PassiveDamage:    passiveDamage,
		IsUpgrade:        isUpgrade,
		AttributableType: dto.Attributable.AttributableName(),
		AttributableID:   dto.Attributable.AttributableID(),
		CreatedAt:        time.Now(), // Лучше использовать текущее время
		UpdatedAt:        time.Now(),
	}).Error

	if err != nil {
		return fmt.Errorf("UserStatService::log: %w", err)
	}

	return nil
}
