package services

import (
	"context"
	"errors"
	"example.com/v2/internal/models"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/errs"
	"example.com/v2/pkg/transaction"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type UserAspectService struct {
	db              *db.DB
	trx             transaction.TransactionManager
	userStatService *UserStatService
}

func NewUserAspectService(
	db *db.DB,
	trx transaction.TransactionManager,
	userStatService *UserStatService,
) *UserAspectService {
	return &UserAspectService{
		db:              db,
		trx:             trx,
		userStatService: userStatService,
	}
}

func (s *UserAspectService) SetAspect(c context.Context, user *models.User, aspect *models.Aspect) error {
	var userAspect models.UserAspect

	err := s.db.WithContext(c).Model(models.UserAspect{}).
		Select("user_aspects.*").
		Joins("inner join aspects a on a.id = user_aspects.aspect_id and a.type = ?", models.Aspects).
		Where("user_id = ?", user.ID).First(&userAspect).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_, err = s.create(c, user, aspect)
			if err != nil {
				return fmt.Errorf("UserAspectService::SetAspect: %w", err)
			}
		} else {
			return fmt.Errorf("UserAspectService::SetAspect: %w", err)
		}
	} else {
		if userAspect.ID == aspect.ID {
			return errs.NewAPIError(
				403,
				fmt.Sprintf("Аспект уже активный"),
			)
		} else if time.Since(userAspect.CreatedAt) < 7*24*time.Hour {
			nextAvailable := userAspect.CreatedAt.Add(7 * 24 * time.Hour)

			return errs.NewAPIError(
				403,
				fmt.Sprintf("Нельзя подключить новый аспект, next_available_date: %s",
					nextAvailable.Format("2006-01-02")),
			)
		}

		err = s.trx.RunInTransaction(c, func(ctx context.Context) error {
			err = s.userStatService.Downgrade(ctx, UserStatUpgradeDto{
				Damage:         userAspect.Damage,
				CriticalDamage: userAspect.CriticalDamage,
				CriticalChance: userAspect.CriticalChance,
				GoldMultiplier: userAspect.GoldMultiplier,
				PassiveDamage:  userAspect.PassiveDamage,
				User:           user,
				Attributable:   aspect,
			})

			if err != nil {
				return fmt.Errorf("UserAspectService::SetAspect: %w", err)
			}

			err = s.db.WithContext(ctx).Model(models.UserAspect{}).Delete(&userAspect).Error

			if err != nil {
				return err
			}

			_, err = s.create(ctx, user, aspect)

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("UserAspectService::SetAspect: %w", err)
		}

		return nil
	}
	return nil
}

func (s *UserAspectService) create(c context.Context, user *models.User, aspect *models.Aspect) (*models.UserAspect, error) {
	var aspectStat models.AspectStat

	err := s.db.WithContext(c).Model(models.AspectStat{}).
		Where("aspect_id = ?", aspect.ID).
		Order("start_level asc").
		First(&aspectStat).Error

	if err != nil {
		return nil, fmt.Errorf("UserAspectService::create: %w", err)
	}

	userAspect := models.UserAspect{
		UserID:         user.ID,
		AspectID:       aspect.ID,
		AspectStatID:   aspectStat.ID,
		Level:          1,
		Damage:         aspectStat.Damage,
		CriticalDamage: aspectStat.CriticalDamage,
		CriticalChance: aspectStat.CriticalChance,
		GoldMultiplier: aspectStat.GoldMultiplier,
		Amount:         aspectStat.Amount,
		PassiveDamage:  aspectStat.PassiveDamage,
	}

	result := s.db.WithContext(c).Create(&userAspect)

	if result.Error != nil {
		return nil, fmt.Errorf("UserAspectService::create: %w", result.Error)
	}

	err = s.userStatService.Upgrade(c, UserStatUpgradeDto{
		Damage:         userAspect.Damage,
		CriticalDamage: userAspect.CriticalDamage,
		CriticalChance: userAspect.CriticalChance,
		GoldMultiplier: userAspect.GoldMultiplier,
		PassiveDamage:  userAspect.PassiveDamage,
		User:           user,
		Attributable:   aspect,
	})

	if err != nil {
		return nil, fmt.Errorf("UserAspectService::create: %w", err)
	}

	return &userAspect, nil
}
