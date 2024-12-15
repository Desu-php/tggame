package services

import (
	"fmt"

	"example.com/v2/internal/models"
	"github.com/gin-gonic/gin"
)

type AuthService struct {}

func NewAuthService(c *gin.Context)*AuthService {
	return &AuthService{}
}

func (s *AuthService) GetUser(c *gin.Context) (*models.User, error){
	userData, exists := c.Get("user")

	if !exists {
		return nil,fmt.Errorf("AuthService::GetUser User not found")
	}

	user, ok := userData.(*models.User)
	if !ok {
		return nil,fmt.Errorf("AuthService::GetUser User data is invalid")
	}

	return user, nil
}