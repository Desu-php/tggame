package responses

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func BadResponse(c *gin.Context, e error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": e.Error(),
	})
}

func OkResponse(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, obj)
}