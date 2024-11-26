package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/v2/internal/adapter"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func SessionMiddleware(sessionCache adapter.UserSessionAdapter, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionHeader := c.GetHeader("x-session")
		if sessionHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing x-session header"})
			c.Abort()
			return
		}

		parts := strings.Split(sessionHeader, "|")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid x-session format"})
			c.Abort()
			return
		}

		telegramID := parts[0]
		session := parts[1]

	    tgId, err := strconv.ParseUint(telegramID, 10, 64)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid x-session format"})
			c.Abort()
			return
		}

		expectedSession, err := sessionCache.Get(tgId)

		if err != nil {
			logger.WithError(err).Error("SessionMiddleware getting session")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server Error"})
			c.Abort()
		}

		if expectedSession != session {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		c.Next()
	}
}