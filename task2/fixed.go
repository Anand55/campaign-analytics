// // main.go (Fixed and Improved Version with Detailed Comments)

package task2

// import (
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"regexp"
// 	"sync"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

// // Global database connection and in-memory campaign spend tracker
// var db *sqlx.DB
// var campaignSpends = make(map[string]float64)
// var mu sync.Mutex // Mutex to ensure safe concurrent map writes

// // Initialize database connection securely using environment variables
// func initDB() {
// 	dsn := os.Getenv("DATABASE_URL")
// 	var err error
// 	db, err = sqlx.Connect("postgres", dsn)
// 	if err != nil {
// 		panic(err) // Fail fast if database cannot be connected
// 	}
// 	fmt.Println("Database connected")
// }

// // Middleware for API Key Authentication
// // Protects all endpoints from unauthorized access
// func authMiddleware(c *gin.Context) {
// 	apiKey := os.Getenv("API_KEY")
// 	reqKey := c.GetHeader("Authorization")
// 	if reqKey != "Bearer "+apiKey {
// 		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}
// 	c.Next()
// }

// // Validate that the campaign ID contains only safe characters
// func validateCampaignID(id string) bool {
// 	re := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
// 	return re.MatchString(id)
// }

// // API Handler to update campaign spend
// func updateSpend(c *gin.Context) {
// 	campaignID := c.Param("campaign_id")
// 	// Validate campaign ID format
// 	if !validateCampaignID(campaignID) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
// 		return
// 	}

// 	// Parse the spend update from request body
// 	var request struct {
// 		Spend float64 `json:"spend"`
// 	}
// 	if err := c.ShouldBindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Safely update in-memory map with lock
// 	mu.Lock()
// 	campaignSpends[campaignID] += request.Spend
// 	mu.Unlock()

// 	// Use a transaction to update the database safely
// 	tx, err := db.Begin()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction start failed"})
// 		return
// 	}

// 	_, err = tx.Exec("UPDATE campaigns SET spend = spend + $1 WHERE id = $2", request.Spend, campaignID)
// 	if err != nil {
// 		tx.Rollback()
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database update failed"})
// 		return
// 	}
// 	tx.Commit()

// 	c.JSON(http.StatusOK, gin.H{"message": "Spend updated"})
// }

// // API Handler to retrieve budget status of a campaign
// func getBudgetStatus(c *gin.Context) {
// 	campaignID := c.Param("campaign_id")
// 	// Validate campaign ID format
// 	if !validateCampaignID(campaignID) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid campaign ID"})
// 		return
// 	}

// 	// Query budget and spend directly from database
// 	var budget, spend float64
// 	err := db.QueryRow("SELECT budget, spend FROM campaigns WHERE id = $1", campaignID).Scan(&budget, &spend)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch campaign data"})
// 		return
// 	}

// 	// Calculate remaining budget and determine status
// 	remaining := budget - spend
// 	status := "Active"
// 	if remaining <= 0 {
// 		status = "Overspent"
// 	}

// 	// Return JSON response with full budget details
// 	c.JSON(http.StatusOK, gin.H{
// 		"campaign_id": campaignID,
// 		"budget":      budget,
// 		"spend":       spend,
// 		"remaining":   remaining,
// 		"status":      status,
// 	})
// }

// // Main function bootstraps the server
// func main() {
// 	initDB()
// 	r := gin.Default()
// 	r.Use(authMiddleware) // Apply authentication middleware globally

// 	// Define routes
// 	r.POST("/campaigns/:campaign_id/spend", updateSpend)
// 	r.GET("/campaigns/:campaign_id/budget-status", getBudgetStatus)

// 	// Start the server on port 8080
// 	r.Run(":8080")
// }
