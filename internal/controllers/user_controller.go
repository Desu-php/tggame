package controllers

import (
	"example.com/v2/internal/repository"
	auth "example.com/v2/internal/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserController struct {
	logger             *logrus.Logger
	auth               *auth.AuthService
	userStatRepository repository.UserStatRepository
	balanceRepository  repository.BalanceRepository
}

func NewUserController(
	logger *logrus.Logger,
	auth *auth.AuthService,
	userStatRepository repository.UserStatRepository,
	balanceRepository repository.BalanceRepository,
) *UserController {
	return &UserController{
		logger:             logger,
		auth:               auth,
		userStatRepository: userStatRepository,
		balanceRepository:  balanceRepository,
	}
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

	c.JSON(http.StatusOK, gin.H{
		"stats":   userStat,
		"balance": balance.Balance,
	})
}
