package controllers

import (
	"errors"
	"example.com/v2/internal/models"
	"example.com/v2/internal/repository"
	"example.com/v2/internal/responses"
	"example.com/v2/internal/services"
	"example.com/v2/pkg/db"
	"example.com/v2/pkg/errs"
	"example.com/v2/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type AspectController struct {
	db                *db.DB
	aspectRepository  repository.AspectRepository
	logger            *logrus.Logger
	userAspectService *services.UserAspectService
}

func NewAspectController(
	db *db.DB,
	aspectRepository repository.AspectRepository,
	logger *logrus.Logger,
	userAspectService *services.UserAspectService,
) *AspectController {
	return &AspectController{
		db:                db,
		aspectRepository:  aspectRepository,
		logger:            logger,
		userAspectService: userAspectService,
	}
}

func (ac *AspectController) Store(c *gin.Context) {
	user, ok := utils.GetUser(c)
	if !ok {
		return
	}

	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	uid, err := strconv.ParseUint(id, 10, 0)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is integer"})
		return
	}

	aspect, err := ac.aspectRepository.FindByID(c, uint(uid), models.Aspects)

	if err == nil && aspect == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aspect not found"})
		return
	}

	if err != nil {
		ac.logger.WithError(err).Error("AspectController::Store")
		responses.ServerErrorResponse(c)
		return
	}

	err = ac.userAspectService.SetAspect(c, user, aspect)

	if err != nil {
		var apiErr *errs.APIError
		if errors.As(err, &apiErr) {
			c.JSON(apiErr.Code, gin.H{"error": apiErr.Error()})
			return
		}
		ac.logger.WithError(err).Error("AspectController::Store")
		responses.ServerErrorResponse(c)
		return
	}
}
