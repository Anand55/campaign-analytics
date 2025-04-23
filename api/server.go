// api/server.go
package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/storage"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for valid API key in Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := os.Getenv("API_KEY")
		header := c.GetHeader("Authorization")
		if expected == "" || !strings.HasPrefix(header, "Bearer ") || strings.TrimPrefix(header, "Bearer ") != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}

// GetCampaignInsights returns the latest metrics for a campaign from cache or DB
func GetCampaignInsights(c *gin.Context) {
	campaignID := c.Param("id")
	from := c.Query("from")
	to := c.Query("to")
	platform := c.Query("platform")

	cacheKey := fmt.Sprintf("campaign:%s:insights:%s:%s:%s", campaignID, from, to, platform)

	cached, err := storage.GetCache(cacheKey)
	if err == nil && cached != "" {
		var response models.CampaignMetrics
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": response, "cached": true})
			return
		}
	}

	query := `SELECT campaign_id, platform, impressions, clicks, conversions, cost, revenue, timestamp
			FROM campaign_metrics WHERE campaign_id = $1`
	args := []interface{}{campaignID}
	argIdx := 2

	if from != "" {
		query += fmt.Sprintf(" AND timestamp >= $%d", argIdx)
		args = append(args, from)
		argIdx++
	}
	if to != "" {
		query += fmt.Sprintf(" AND timestamp <= $%d", argIdx)
		args = append(args, to)
		argIdx++
	}
	if platform != "" {
		query += fmt.Sprintf(" AND platform = $%d", argIdx)
		args = append(args, platform)
		argIdx++
	}

	query += " ORDER BY timestamp DESC LIMIT 1"

	row := storage.DB.QueryRow(query, args...)
	var result models.CampaignMetrics
	err = row.Scan(
		&result.CampaignID,
		&result.Platform,
		&result.Impressions,
		&result.Clicks,
		&result.Conversions,
		&result.Cost,
		&result.Revenue,
		&result.Timestamp,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found for campaign"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query DB"})
		return
	}

	serialized, _ := json.Marshal(result)
	storage.SetCache(cacheKey, string(serialized), 30*time.Second)

	c.JSON(http.StatusOK, gin.H{"data": result, "cached": false})
}

// InitRouter sets up the Gin router and routes
func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(AuthMiddleware())
	r.GET("/campaign/:id/insights", GetCampaignInsights)

	return r
}
