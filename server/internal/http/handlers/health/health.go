package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter(router *gin.RouterGroup) {
	health := router.Group("/health")
	{
		health.GET("", healthCheck)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
