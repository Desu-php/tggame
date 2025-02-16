package controllers

import (
	"example.com/v2/internal/http/resources"
	"example.com/v2/pkg/image"
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
	image   *image.Image
}

func NewGameController(service *services.GameService, logger *logrus.Logger, image *image.Image) *GameController {
	return &GameController{
		service: service,
		logger:  logger,
		image:   image,
	}
}

func (gc *GameController) Start(c *gin.Context) {
	var startDto dto.GameStartDto

	if err := c.ShouldBindJSON(&startDto); err != nil {
		responses.BadResponse(c, err)
		return
	}

	user, err := gc.service.Start(c, &startDto)

	if err != nil {
		gc.logger.WithContext(c).Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Server error",
		})
		return
	}

	responses.OkResponse(c, gin.H{"user": resources.NewUserResource(gc.image).Map(user)})
}
