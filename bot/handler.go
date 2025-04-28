package bot

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler for receiving user prompts
func PromptHandler(c *gin.Context) {
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Step 1: Embed the user's prompt into a vector
	embedding, err := EmbedText(req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to embed prompt"})
		return
	}

	// Step 2: Search nearest campaign using vector similarity
	campaignID, err := SearchNearestCampaign(embedding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No matching campaign found"})
		return
	}

	// Step 3: Query your Analytics API to fetch campaign insights
	data, err := QueryAnalyticsBackend(campaignID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
		return
	}

	// Step 4: Format and summarize the response
	summary := FormatResponse("summary", data) // ✅ "summary" passed as first argument

	// Step 5: Return the final response to user
	c.JSON(http.StatusOK, gin.H{"response": summary})
}
