package controllers

import (
	"net/http"

	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/dto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GameController struct {
	service *services.GameService
	logger  *logrus.Logger
}

func NewGameController(service *services.GameService, logger *logrus.Logger) *GameController {
	return &GameController{
		service: service,
		logger:  logger,
	}
}

type StartResponse struct {
	ID         uint               `json:"id"`
	Session    string             `json:"session"`
	TelegramID uint64             `json:"telegram_id"`
	UserChest  *UserChestResponse `json:"user_chest"`
}

type UserChestResponse struct {
	Health        uint   `json:"health"`
	CurrentHealth int    `json:"current_health"`
	Level         int    `json:"level"`
	Name          string `json:"name"`
	Image         string `json:"url"`
}

func (gc *GameController) Start(c *gin.Context) {
	var dto dto.GameStartDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		responses.BadResponse(c, err)
		return
	}

	user, err := gc.service.Start(&dto)

	if err != nil {
		gc.logger.WithContext(c).Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error",
		})
		return
	}

	responses.OkResponse(c, gin.H{"user": user})
}
