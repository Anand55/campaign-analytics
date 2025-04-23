// api/server.go
package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/storage"

	"github.com/gin-gonic/gin"
)

// GetCampaignInsights returns the latest metrics for a campaign from cache or DB
func GetCampaignInsights(c *gin.Context) {
	campaignID := c.Param("id")
	cacheKey := fmt.Sprintf("campaign:%s:insights", campaignID)

	// 1. Try Redis Cache
	cached, err := storage.GetCache(cacheKey)
	if err == nil && cached != "" {
		var response models.CampaignMetrics
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			c.JSON(http.StatusOK, gin.H{"data": response, "cached": true})
			return
		}
	}

	// 2. Query Postgres for the latest campaign metrics
	query := `SELECT campaign_id, platform, impressions, clicks, conversions, cost, revenue, timestamp
			FROM campaign_metrics
			WHERE campaign_id = $1
			ORDER BY timestamp DESC LIMIT 1`

	row := storage.DB.QueryRow(query, campaignID)
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

	// 3. Cache the result
	serialized, _ := json.Marshal(result)
	storage.SetCache(cacheKey, string(serialized), 30*time.Second)

	// 4. Return the response
	c.JSON(http.StatusOK, gin.H{"data": result, "cached": false})
}

// InitRouter sets up the Gin router and routes
func InitRouter() *gin.Engine {
	r := gin.Default()

	// Register routes
	r.GET("/campaign/:id/insights", GetCampaignInsights)

	return r
}
