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
}

func NewUserController(
	logger *logrus.Logger,
	auth *auth.AuthService,
	userStatRepository repository.UserStatRepository,
) *UserController {
	return &UserController{
		logger:             logger,
		auth:               auth,
		userStatRepository: userStatRepository,
	}
}

func (uc *UserController) Info(c *gin.Context) {
	user, err := uc.auth.GetUser(c)

	if err != nil {
		uc.logger.WithError(err).Error("UserController::Info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	userStat, err := uc.userStatRepository.GetStat(user)

	if err != nil {
		uc.logger.WithError(err).Error("UserController::Info")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": userStat,
	})
}
