package controllers

import (
	"example.com/v2/internal/repository"
	"example.com/v2/internal/responses"
	auth "example.com/v2/internal/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserItemController struct {
	userItemRepository repository.UserItemRespository
	logger             *logrus.Logger
	auth               *auth.AuthService
}

func NewUserItemController(userItemRepository repository.UserItemRespository, logger *logrus.Logger) *UserItemController {
	return &UserItemController{
		userItemRepository: userItemRepository,
		logger:             logger,
	}
}

func (cc *UserItemController) GetLast(c *gin.Context) {
	user, err := cc.auth.GetUser(c)

	if err != nil {
		cc.logger.WithError(err).Error("UserItemController::GetLast")
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}

	userItem, err := cc.userItemRepository.GetLast(user.ID)

	if err != nil {
		cc.logger.WithError(err).Error("UserItemController::GetLast")
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}

	responses.OkResponse(c, gin.H{"item": userItem})
}


func (cc *UserItemController) GetUserItems(c *gin.Context) {
	user, err := cc.auth.GetUser(c)

	if err != nil {
		cc.logger.WithError(err).Error("UserItemController::GetLast")
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}

	userItems, err := cc.userItemRepository.GetUserItems(user.ID)

	if err != nil {
		cc.logger.WithError(err).Error("UserItemController::GetLast")
		c.JSON(500, gin.H{"error": "Server error"})
		return
	}

	responses.OkResponse(c, userItems)
}