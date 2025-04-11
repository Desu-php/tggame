package controllers

import (
	"context"
	"errors"
	"example.com/v2/internal/http/resources"
	responses2 "example.com/v2/internal/http/responses"
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/errs"
	"example.com/v2/pkg/image"
	"example.com/v2/pkg/transaction"
	"example.com/v2/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type BoosterController struct {
	logger          *logrus.Logger
	image           *image.Image
	db              *db.DB
	trx             transaction.TransactionManager
	userStatService *services.UserStatService
	balanceService  *services.BalanceService
}

func NewBoosterController(
	logger *logrus.Logger,
	image *image.Image,
	db *db.DB,
	trx transaction.TransactionManager,
	userStatService *services.UserStatService,
	balanceService *services.BalanceService,
) *BoosterController {
	return &BoosterController{
		logger,
		image,
		db,
		trx,
		userStatService,
		balanceService,
	}
}

func (as *BoosterController) Index(c *gin.Context) {
	user, ok := utils.GetUser(c)
	if !ok {
		return
	}

	as.db.WithContext(c).Model(models.UserAspect{}).First(&user)

	var aspects []responses2.AspectWithStatsResponse

	err := as.db.WithContext(c).
		Raw(`
SELECT
  a.id,
  a.name,
  a.image,
  a.description,
  ua.level AS user_level,
  ast.id as aspect_stat_id,
  ast.damage,
  ast.critical_damage,
  ast.critical_chance,
  ast.gold_multiplier,
  ast.amount,
  ast.amount_growth_factor
FROM aspects a
LEFT JOIN user_aspects ua ON ua.aspect_id = a.id AND ua.user_id = ?
LEFT JOIN LATERAL (
  SELECT *
  FROM aspect_stats ast
  WHERE ast.aspect_id = a.id
    AND (
      (ua.level BETWEEN ast.start_level AND ast.end_level)
      OR ua.level IS NULL
    )
  ORDER BY
    CASE
      WHEN ua.level BETWEEN ast.start_level AND ast.end_level THEN 0
      ELSE 1
    END,
    start_level
  LIMIT 1
) ast ON true
WHERE a.type = ?
`, user.ID, models.Booster).Scan(&aspects).Error

	if err != nil {
		as.logger.WithError(err).Error("BoosterController::index")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{
		"data": resources.NewBaseResource(resources.NewAspectWithStatsResource(as.image)).All(aspects),
	})
}

func (as *BoosterController) Buy(c *gin.Context) {
	user, ok := utils.GetUser(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var aspect models.Aspect
	aspectType := models.Booster

	err := as.db.WithContext(c).
		Where("type = ?", aspectType).
		First(&aspect, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booster not found"})
			return
		}
		as.logger.WithError(err).Error("BoosterController::buy")
		responses.ServerErrorResponse(c)
		return
	}

	var userAspect models.UserAspect
	err = as.db.WithContext(c).
		Where("user_id = ? AND aspect_id = ?", user.ID, aspect.ID).
		First(&userAspect).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already owns this booster",
			"code":  "booster_already_owned",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		as.logger.WithError(err).Error("BoosterController::buy")
		responses.ServerErrorResponse(c)
		return
	}

	var aspectStat models.AspectStat

	err = as.db.WithContext(c).
		Order("start_level asc").
		First(&aspectStat, "aspect_id = ?", aspect.ID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Booster stats not found",
				"details": fmt.Sprintf("No stats found for booster ID: %s", aspect.ID),
			})
			return
		}
		as.logger.WithError(err).Error("BoosterController::buy")
		responses.ServerErrorResponse(c)
		return
	}

	err = as.trx.RunInTransaction(c, func(ctx context.Context) error {

		newUserAspect := models.UserAspect{
			UserID:         user.ID,
			AspectID:       aspect.ID,
			AspectStatID:   aspectStat.ID,
			Level:          1,
			Damage:         aspectStat.Damage,
			CriticalDamage: aspectStat.CriticalDamage,
			CriticalChance: aspectStat.CriticalChance,
			GoldMultiplier: aspectStat.GoldMultiplier,
			Amount:         aspectStat.Amount,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := as.db.WithContext(ctx).Create(&newUserAspect).Error; err != nil {
			return fmt.Errorf("failed to create user aspect: %w", err)
		}

		err = as.balanceService.Charge(ctx, &services.TransactionDto{
			Amount: int64(aspectStat.Amount),
			User:   user,
			Model:  &aspectStat,
			Type:   models.TransactionTypeBuyAspect,
		})

		if err != nil {
			return err
		}

		err := as.userStatService.Upgrade(ctx, services.UserStatUpgradeDto{
			Damage:         newUserAspect.Damage,
			CriticalDamage: newUserAspect.CriticalDamage,
			CriticalChance: newUserAspect.CriticalChance,
			GoldMultiplier: newUserAspect.GoldMultiplier,
			User:           user,
		})

		if err != nil {
			return fmt.Errorf("failed to upgrade user stat: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, errs.ErrInsufficientBalance) {
			c.JSON(400, gin.H{"error": "Недостаточно средств на балансе"})
			return
		}

		as.logger.WithError(err).Error("BoosterController::buy")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{})
}

func (as *BoosterController) Upgrade(c *gin.Context) {
	user, ok := utils.GetUser(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var aspect models.Aspect
	aspectType := models.Booster

	err := as.db.WithContext(c).
		Where("type = ?", aspectType).
		First(&aspect, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booster not found"})
			return
		}
		as.logger.WithError(err).Error("BoosterController::Upgrade")
		responses.ServerErrorResponse(c)
		return
	}

	var userAspect models.UserAspect
	err = as.db.WithContext(c).
		Where("user_id = ? AND aspect_id = ?", user.ID, aspect.ID).
		First(&userAspect).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User aspect not found"})
			return
		}

		as.logger.WithError(err).Error("BoosterController::Upgrade")
		responses.ServerErrorResponse(c)
		return
	}

	userAspect.Level = userAspect.Level + 1

	var aspectStat models.AspectStat

	err = as.db.WithContext(c).
		Where("? BETWEEN start_level AND end_level", userAspect.Level).
		First(&aspectStat, "aspect_id = ?", aspect.ID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(400, gin.H{"error": "Бустер достиг максимального уровня"})
			return
		} else {
			as.logger.WithError(err).Error("BoosterController::Upgrade")
			responses.ServerErrorResponse(c)
			return
		}
	}

	err = as.trx.RunInTransaction(c, func(ctx context.Context) error {
		var amount int64

		log.Println(aspectStat)

		if aspectStat.ID != userAspect.AspectStatID {
			amount = int64(aspectStat.Amount)
		} else {
			amount = int64(utils.GrowthIncrease(float64(userAspect.Amount), aspectStat.AmountGrowthFactor))
		}

		userAspect.AspectStatID = aspectStat.ID
		userAspect.CriticalDamage += aspectStat.CriticalDamage
		userAspect.Damage += aspectStat.Damage
		userAspect.CriticalChance += aspectStat.CriticalChance
		userAspect.GoldMultiplier += aspectStat.GoldMultiplier
		userAspect.Amount = uint(amount)

		if err := as.db.WithContext(ctx).Model(models.UserAspect{}).Where("id = ?", userAspect.ID).Updates(userAspect).Error; err != nil {
			return fmt.Errorf("failed to update user aspect: %w", err)
		}

		err = as.balanceService.Charge(ctx, &services.TransactionDto{
			Amount: amount,
			User:   user,
			Model:  &aspectStat,
			Type:   models.TransactionTypeBuyAspect,
		})

		if err != nil {
			return err
		}

		err := as.userStatService.Upgrade(ctx, services.UserStatUpgradeDto{
			Damage:         aspectStat.Damage,
			CriticalDamage: aspectStat.CriticalDamage,
			CriticalChance: aspectStat.CriticalChance,
			GoldMultiplier: aspectStat.GoldMultiplier,
			User:           user,
		})

		if err != nil {
			return fmt.Errorf("failed to upgrade user stat: %w", err)
		}

		return nil
	})

	if err != nil {
		if errors.Is(err, errs.ErrInsufficientBalance) {
			c.JSON(400, gin.H{"error": "Недостаточно средств на балансе"})
			return
		}

		as.logger.WithError(err).Error("BoosterController::Upgrade")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{})
}
