package webhook

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	wg.Add(3)

	// Bonus B: Concurrent tasks using Goroutines and sync.WaitGroup
	go writeLog(&wg, payload)
	go updateDbStats(&wg, payload)
	go asyncNotification(&wg, payload)

	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"status": "Webhook processed concurrently"})
}

func writeLog(wg *sync.WaitGroup, payload map[string]interface{}) {
	defer wg.Done()
	f, err := os.OpenFile("system.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Could not write to log file: ", err)
		return
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("Received: %v\n", payload))
}

func updateDbStats(wg *sync.WaitGroup, payload map[string]interface{}) {
	defer wg.Done()
	// Mock DB update
	log.Println("[Concurrency] DB Stats updated for:", payload["type"])
}

func asyncNotification(wg *sync.WaitGroup, payload map[string]interface{}) {
	defer wg.Done()
	// Mock Notification
	log.Println("[Concurrency] Notification dispatched.")
}
