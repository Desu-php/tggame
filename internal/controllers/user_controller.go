package controllers

import (
	"errors"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/internal/responses"
	auth "example.com/v2/internal/services/auth"
	"example.com/v2/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	logger             *logrus.Logger
	auth               *auth.AuthService
	userStatRepository repository.UserStatRepository
	balanceRepository  repository.BalanceRepository
	userRepository     repository.UserRepository
	db                 *db.DB
}

func NewUserController(
	logger *logrus.Logger,
	auth *auth.AuthService,
	userStatRepository repository.UserStatRepository,
	balanceRepository repository.BalanceRepository,
	userRepository repository.UserRepository,
	db *db.DB,
) *UserController {
	return &UserController{
		logger:             logger,
		auth:               auth,
		userStatRepository: userStatRepository,
		balanceRepository:  balanceRepository,
		userRepository:     userRepository,
		db:                 db,
	}
}

type GroupedItem struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ItemsCount uint   `json:"items_count"`
}

func (uc *UserController) Info(c *gin.Context) {
	user, err := uc.auth.GetUser(c)

	if err != nil {
		uc.logger.WithError(err).Error("UserController::Info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	userStat, err := uc.userStatRepository.GetStat(c, user)

	if err != nil {
		uc.logger.WithError(err).Error("UserController::Info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	balance, err := uc.balanceRepository.FindByUserId(c, user.ID)

	if err != nil {
		uc.logger.WithError(err).Error("UserController::Info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	var referralUser models.User

	err = uc.db.WithContext(c).Model(models.User{}).
		Select("users.*").
		Joins("inner join referral_users ru ON ru.user_id = users.id and ru.referred_user_id = ?", user.ID).
		First(&referralUser).Error

	var invite *string = nil

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == false {
			uc.logger.WithError(err).Error("UserController::Info")
			responses.ServerErrorResponse(c)
			return
		}
	} else {
		invite = &referralUser.Username
	}

	var items []GroupedItem

	err = uc.db.WithContext(c).Model(models.UserItem{}).
		Select("r.id, r.name, count(DISTINCT item_id) as items_count").
		Joins("inner join items as i on i.id = user_items.item_id").
		Joins("inner join rarities as r on r.id = i.rarity_id").
		Group("r.id").
		Order("r.sort asc").
		Find(&items).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) == false {
			uc.logger.WithError(err).Error("UserController::Info")
			responses.ServerErrorResponse(c)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats":   userStat,
		"balance": balance.Balance,
		"invite":  invite,
		"items":   items,
	})
}

func (uc *UserController) GetTop(c *gin.Context) {
	users, err := uc.userRepository.GetTop(c)

	if err != nil {
		uc.logger.WithError(err).Error("UserController::GetTop")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}
