package controllers

import (
	"example.com/v2/internal/http/resources"
	"example.com/v2/internal/models"
	"example.com/v2/internal/responses"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/image"
	"example.com/v2/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserAspectController struct {
	db     *db.DB
	image  *image.Image
	logger *logrus.Logger
}

func NewUserAspectController(db *db.DB, image *image.Image, logger *logrus.Logger) *UserAspectController {
	return &UserAspectController{db, image, logger}
}

func (a *UserAspectController) GetBoosters(c *gin.Context) {
	user, ok := utils.GetUser(c)

	if !ok {
		return
	}

	var aspects []models.Aspect

	err := a.db.WithContext(c).Model(models.Aspect{}).
		Select("aspects.*").
		Joins("inner join user_aspects ua ON ua.aspect_id = aspects.id").
		Where("ua.user_id = ?", user.ID).
		Where("aspects.type = ?", models.Booster).
		Find(&aspects).Error

	if err != nil {
		a.logger.WithError(err).Error("UserAspectController::GetBoosters failed")
	}

	responses.OkResponse(c, gin.H{
		"data": resources.NewBaseResource(resources.NewAspectResource(a.image)).All(aspects),
	})
}
