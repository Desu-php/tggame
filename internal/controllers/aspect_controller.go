package controllers

import (
	"context"
	"errors"
	"example.com/v2/internal/http/resources"
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/image"
	"example.com/v2/pkg/transaction"
	"example.com/v2/pkg/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type AspectController struct {
	logger          *logrus.Logger
	image           *image.Image
	db              *db.DB
	trx             transaction.TransactionManager
	userStatService *services.UserStatService
	balanceService  *services.BalanceService
}

func NewAspectController(
	logger *logrus.Logger,
	image *image.Image,
	db *db.DB,
	trx transaction.TransactionManager,
	userStatService *services.UserStatService,
	balanceService *services.BalanceService,
) *AspectController {
	return &AspectController{
		logger,
		image,
		db,
		trx,
		userStatService,
		balanceService,
	}
}

func (as *AspectController) Index(c *gin.Context) {
	var aspects []models.Aspect

	err := as.db.WithContext(c).
		Model(models.Aspect{}).
		Find(&aspects).Error

	if err != nil {
		as.logger.WithError(err).Error("AspectController::index")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{
		"data": resources.NewBaseResource(resources.NewAspectResource(as.image)).All(aspects),
	})
}

func (as *AspectController) Buy(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	user, ok := userData.(*models.User)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user data"})
		return
	}

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var aspect models.Aspect
	aspectType := models.Aspects

	err := as.db.WithContext(c).
		Where("type = ?", aspectType).
		First(&aspect, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aspect not found"})
			return
		}
		as.logger.WithError(err).Error("AspectController::buy")
		responses.ServerErrorResponse(c)
		return
	}

	var userAspect models.UserAspect
	err = as.db.WithContext(c).
		Where("user_id = ? AND aspect_id = ?", user.ID, aspect.ID).
		First(&userAspect).Error

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "User already owns this aspect",
			"code":  "aspect_already_owned",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		as.logger.WithError(err).Error("AspectController::buy")
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
				"error":   "Aspect stats not found",
				"details": fmt.Sprintf("No stats found for aspect ID: %s", aspect.ID),
			})
			return
		}
		as.logger.WithError(err).Error("AspectController::buy")
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
			return fmt.Errorf("failed to charge user stat: %w", err)
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
		as.logger.WithError(err).Error("AspectController::buy")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{})
}

func (as *AspectController) Upgrade(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	user, ok := userData.(*models.User)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user data"})
		return
	}

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var aspect models.Aspect
	aspectType := models.Aspects

	err := as.db.WithContext(c).
		Where("type = ?", aspectType).
		First(&aspect, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Aspect not found"})
			return
		}
		as.logger.WithError(err).Error("AspectController::Upgrade")
		responses.ServerErrorResponse(c)
		return
	}

	var userAspect models.UserAspect
	err = as.db.WithContext(c).
		Where("user_id = ? AND aspect_id = ?", user.ID, aspect.ID).
		First(&userAspect).Error

	if err != nil {
		as.logger.WithError(err).Error("AspectController::Upgrade")
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
			err = as.db.WithContext(c).
				Order("end_level desc").
				Where("end_level <= ?", userAspect.Level).
				First(&aspectStat, "aspect_id = ?", aspect.ID).Error

			if err != nil {
				as.logger.WithError(err).Error("AspectController::Upgrade")
				responses.ServerErrorResponse(c)
				return
			}
		} else {
			as.logger.WithError(err).Error("AspectController::Upgrade")
			responses.ServerErrorResponse(c)
			return
		}
	}

	err = as.trx.RunInTransaction(c, func(ctx context.Context) error {
		var amount int64

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
			return fmt.Errorf("failed to charge user stat: %w", err)
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
		as.logger.WithError(err).Error("AspectController::Upgrade")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{})
}
