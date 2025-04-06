package utils

import (
	"example.com/v2/internal/models"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

func GrowthIncrease(currentValue float64, growthFactor float64) float64 {
	increase := currentValue * (growthFactor / 100)
	return math.Round(currentValue + increase)
}

func GetUser(c *gin.Context) (*models.User, bool) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return nil, false
	}

	user, ok := userData.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return nil, false
	}

	return user, true
}
