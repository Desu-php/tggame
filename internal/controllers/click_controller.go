package controllers

import (
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClickStoreRequest struct {
	Count  uint `json:"count" binding:"required"`
}

type ClickController struct {
	logger  *logrus.Logger
	service *services.ClickService
}

func NewClickController(
	logger *logrus.Logger,
	service *services.ClickService,
) *ClickController {
	return &ClickController{
		logger:  logger,
		service: service,
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
		c.JSON(401, gin.H{"error": "User not found"})
		return
	}

	user, ok := userData.(*models.User)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user data"})
		return
	}

	cc.service.Damage(user,request.Count)

	responses.OkResponse(c, gin.H{"user": "Hello world!"})
}
