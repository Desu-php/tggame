package controllers

import (
	"net/http"

	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/dto"
	"github.com/gin-gonic/gin"
)

type GameController struct {
	service *services.GameService
}

func NewGameController(service *services.GameService) *GameController {
	return &GameController{service: service}
}

func (gc *GameController) Start(c *gin.Context) {

	var dto dto.GameStartDto

	if err := c.ShouldBindJSON(&dto); err != nil {
		responses.BadResponse(c, err)
		return
	}

	user, err := gc.service.Start(&dto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Something went wrong",
		})
		return
	}

	responses.OkResponse(c, gin.H{"user": user})
}
