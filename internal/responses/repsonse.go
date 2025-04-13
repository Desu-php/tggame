package responses

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BadResponse(c *gin.Context, e error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"error": e.Error(),
	})
}

func OkResponse(c *gin.Context, obj any) {
	c.JSON(http.StatusOK, obj)
}

func ServerErrorResponse(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, "Server Error")
}

func ServerErrorResponseWithMessage(c *gin.Context, obj any) {
	c.JSON(http.StatusInternalServerError, obj)
}

func NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, "Not Found")
}
