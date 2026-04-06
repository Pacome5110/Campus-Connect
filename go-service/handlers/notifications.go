package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostNotification(c *gin.Context) {
	var payload struct {
		Message string `json:"message"`
		UserId  string `json:"userId"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Notification sent successfully",
		"data":   payload,
	})
}
