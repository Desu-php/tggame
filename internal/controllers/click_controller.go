package controllers

import (
	"log"

	"example.com/v2/internal/adapter"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ClickStoreRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Count  uint `json:"count" binding:"required"`
}

type ClickController struct {
	logger  *logrus.Logger
	service *services.ClickService
	cache   adapter.UserClickCacheAdapter
}

func NewClickController(
	logger *logrus.Logger,
	service *services.ClickService,
	cache adapter.UserClickCacheAdapter,
) *ClickController {
	return &ClickController{
		logger:  logger,
		service: service,
		cache:   cache,
	}
}

func (cc *ClickController) Store(c *gin.Context) {
	var request ClickStoreRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		responses.BadResponse(c, err)
		return
	}

	cc.service.Store(request.UserID, request.Count)

	responses.OkResponse(c, gin.H{"user": "Hello world!"})
}

func (cc *ClickController) Get(c *gin.Context) {

	cc.cache.ChunkAll(5, func (data map[uint]uint) error  {
		log.Println("len", len(data))
		for userId, value := range data {
			log.Println("value: ", userId, value)
		}
		return nil
	})
		responses.OkResponse(c, gin.H{"user": "Hello world!"})
}
