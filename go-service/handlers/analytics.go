package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAnalytics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"event_count": 142,
		"active_users": 530,
	})
}
