package controllers

import (
	"example.com/v2/internal/http/resources"
	"example.com/v2/internal/repository"
	auth "example.com/v2/internal/services/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type ReferralController struct {
	logger     *logrus.Logger
	repository repository.ReferralUserRepository
	auth       *auth.AuthService
}

func NewReferralController(
	logger *logrus.Logger,
	repository repository.ReferralUserRepository,
	auth *auth.AuthService,
) *ReferralController {
	return &ReferralController{
		logger:     logger,
		repository: repository,
		auth:       auth,
	}
}

func (rc *ReferralController) GetReferrals(c *gin.Context) {
	user, err := rc.auth.GetUser(c)

	if err != nil {
		rc.logger.WithError(err).Error("NewReferralController::getReferrals")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	referrals, err := rc.repository.GetByUserID(user.ID)

	log.Println(referrals)

	if err != nil {
		rc.logger.WithError(err).Error("NewReferralController::getReferrals")
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": resources.NewBaseResource(resources.NewReferralUserResource()).All(referrals),
		},
	)
}

func (rc *ReferralController) GetReferralCount(c *gin.Context) {
	user, err := rc.auth.GetUser(c)

	if err != nil {
		rc.logger.WithError(err).Error("NewReferralController::GetReferralCount")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	count, err := rc.repository.Count(user.ID)

	if err != nil {
		rc.logger.WithError(err).Error("NewReferralController::GetReferralCount")
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
