package middleware

import (
	"strings"
	"example.com/v2/config"
	"example.com/v2/pkg/telegram"
	"github.com/gin-gonic/gin"
)

func TelegramHashCheck(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {

		if !cfg.IsProduction() {
			c.Next()
			return
		}

		if !isValid(cfg.Telegram.Token, c.GetHeader("X-Telegram-Init-Data")) {
			c.JSON(401, gin.H{"error": "Invalid initData"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isValid(botTokens, initData string) bool {
	tokens := strings.Split(botTokens, ",")

	for _, token := range tokens {
		if telegram.ValidateInitData(initData, token) {
			return true
		}
	}

	return false
}
