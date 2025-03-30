package controllers

import (
	"example.com/v2/internal/http/resources"
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/image"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AspectController struct {
	logger *logrus.Logger
	image  *image.Image
	db     *db.DB
}

func NewAspectController(logger *logrus.Logger, image *image.Image, db *db.DB) *AspectController {
	return &AspectController{logger, image, db}
}

func (as *AspectController) Index(c *gin.Context) {
	var aspects []models.Aspect

	err := as.db.WithContext(c).
		Model(models.Aspect{}).
		Find(&aspects).Error

	if err != nil {
		as.logger.WithError(err).Error("AspectController::index")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{
		"data": resources.NewBaseResource(resources.NewAspectResource(as.image)).All(aspects),
	})
}
