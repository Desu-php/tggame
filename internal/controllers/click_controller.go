package controllers

import (
	"example.com/v2/internal/http/resources"
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/errs"
	"example.com/v2/pkg/image"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClickStoreRequest struct {
	Count uint `json:"count" binding:"required"`
}

type ClickController struct {
	logger  *logrus.Logger
	service *services.ClickService
	image   *image.Image
}

func NewClickController(
	logger *logrus.Logger,
	service *services.ClickService,
	image *image.Image,
) *ClickController {
	return &ClickController{
		logger:  logger,
		service: service,
		image:   image,
	}
}

func (cc *ClickController) Store(c *gin.Context) {
	var request ClickStoreRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		responses.BadResponse(c, err)
		return
	}

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(401, gin.H{"error": "User not found", "code": errs.UnauthorizedCode})
		return
	}

	user, ok := userData.(*models.User)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user data"})
		return
	}

	err := cc.service.Damage(c, user, request.Count)

	if err != nil {
		cc.logger.WithError(err).Error("ClickController::store")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{
		"data": resources.NewBaseResource(resources.NewUserChestResource(cc.image)).One(&user.UserChest),
	})
}
