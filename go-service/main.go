package main

import (
	"log"

	"campusconnect-go/handlers"
	"campusconnect-go/middleware"
	"campusconnect-go/webhook"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rate Limiter
	r.Use(middleware.RateLimiter(5))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/webhook", webhook.HandleWebhook)

	api := r.Group("/api")
	api.Use(middleware.APIKeyAuth())
	{
		api.GET("/analytics", handlers.GetAnalytics)
		api.POST("/notifications", handlers.PostNotification)
	}

	return r
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, defaulting to system enc")
	}

	r := SetupRouter()
	log.Println("Go service starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
