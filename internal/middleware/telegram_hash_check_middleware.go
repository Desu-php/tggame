package middleware

import (
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

		if !telegram.ValidateInitData(c.GetHeader("X-Telegram-Init-Data"), cfg.Telegram.Token) {
			c.JSON(401, gin.H{"error": "Invalid initData"})
			c.Abort()
			return
		}
		
		c.Next()
	}
}