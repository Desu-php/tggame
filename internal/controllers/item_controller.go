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

type ItemController struct {
	db     *db.DB
	logger *logrus.Logger
	image  *image.Image
}

func NewItemController(db *db.DB, logger *logrus.Logger, image *image.Image) *ItemController {
	return &ItemController{db, logger, image}
}

func (cc *ItemController) GetItems(c *gin.Context) {
	var items []models.Item

	query := cc.db.WithContext(c).Model(models.Item{}).
		Preload("Rarity").
		Preload("Type")

	rarityId := c.Query("rarity_id")
	if rarityId != "" {
		query = query.Where("rarity_id = ?", rarityId)
	}

	err := query.Find(&items).Error

	if err != nil {
		cc.logger.WithError(err).Error("ItemController::GetItems")
		responses.ServerErrorResponse(c)
		return
	}

	responses.OkResponse(c, gin.H{
		"data": resources.NewBaseResource(resources.NewItemResource(cc.image)).All(items),
	})
}
