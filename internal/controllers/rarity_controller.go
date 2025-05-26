package controllers

import (
	"example.com/v2/internal/http/resources"
	repository "example.com/v2/internal/repository/item"
	"example.com/v2/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RarityController struct {
	logger     *logrus.Logger
	repository repository.RarityRepository
}

func NewRarityController(logger *logrus.Logger, repository repository.RarityRepository) *RarityController {
	return &RarityController{
		logger:     logger,
		repository: repository,
	}
}

func (cc *RarityController) GetRarities(c *gin.Context) {
	rarities, err := cc.repository.GetAll(c)

	if err != nil {
		cc.logger.WithContext(c).Errorf("failed to get rarities: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	responses.OkResponse(c, resources.NewBaseResource(resources.NewRarityResource()).All(rarities))
}
