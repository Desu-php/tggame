package middleware

import (
	"example.com/v2/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogHandle(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				logger.WithError(e).Error(e.Error())
			}

			if c.Writer.Status() < 400 {
				responses.ServerErrorResponse(c)
			}
		}
	}
}
