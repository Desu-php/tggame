package controllers

import (
	"example.com/v2/internal/http/resources"
	"example.com/v2/internal/repository"
	"example.com/v2/internal/responses"
	auth "example.com/v2/internal/services/auth"
	"example.com/v2/pkg/image"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserItemController struct {
	userItemRepository repository.UserItemRespository
	logger             *logrus.Logger
	auth               *auth.AuthService
	image              *image.Image
}

func NewUserItemController(userItemRepository repository.UserItemRespository, logger *logrus.Logger, image *image.Image) *UserItemController {
	return &UserItemController{
		userItemRepository: userItemRepository,
		logger:             logger,
		image:              image,
	}
}

func (cc *UserItemController) GetLast(c *gin.Context) {
	user, err := cc.auth.GetUser(c)

	if err != nil {
		cc.logger.WithError(err).Error("UserItemController::GetLast")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	userItem, err := cc.userItemRepository.GetLast(user.ID)

	if err != nil {
		cc.logger.WithError(err).Error("UserItemController::GetLast")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	responses.OkResponse(c, gin.H{"data": resources.NewUserItemResource(cc.image).Map(userItem)})
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

	responses.OkResponse(c,
		resources.NewBaseResource(resources.NewGroupedUserItemResource(cc.image)).All(userItems),
	)
}
