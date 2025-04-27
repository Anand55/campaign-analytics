// ingestion/linkedin.go
package ingestion

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/processor"
)

// FetchLinkedInCampaigns pulls campaign insights from LinkedIn Marketing API
func FetchLinkedInCampaigns() {
	fmt.Println("[LINKEDIN] Fetching campaign data from LinkedIn Marketing API...")

	token := os.Getenv("LINKEDIN_ACCESS_TOKEN")
	accountID := os.Getenv("LINKEDIN_ACCOUNT_ID")

	if token == "" || accountID == "" {
		fmt.Println("[LINKEDIN] Missing LinkedIn API credentials.")
		return
	}

	// start := time.Now().AddDate(0, 0, -7).Unix() * 1000 // 7 days ago in ms
	// end := time.Now().Unix() * 1000                     // now in ms

	url := fmt.Sprintf("https://api.linkedin.com/v2/adAnalyticsV2?q=analytics&dateRange.start.day=1&dateRange.start.month=1&dateRange.start.year=2024&dateRange.end.day=30&dateRange.end.month=4&dateRange.end.year=2024&pivot=CAMPAIGN&accounts=urn:li:sponsoredAccount:%s", accountID)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[LINKEDIN] Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[LINKEDIN] API returned non-200: %d\n", resp.StatusCode)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var response struct {
		elements []struct {
			CampaignURN         string  `json:"campaign"`
			Impressions         int     `json:"impressions"`
			Clicks              int     `json:"clicks"`
			CostInLocalCurrency float64 `json:"costInLocalCurrency"`
		} `json:"elements"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Printf("[LINKEDIN] Failed to parse response: %v\n", err)
		return
	}

	for _, item := range response.elements {
		// Campaign ID from URN: urn:li:sponsoredCampaign:<campaign_id>
		parts := strings.Split(item.CampaignURN, ":")
		campaignID := "l-unknown"
		if len(parts) == 4 {
			campaignID = "l-" + parts[3]
		}

		metric := models.CampaignMetrics{
			CampaignID:  campaignID,
			Platform:    "LinkedIn",
			Impressions: item.Impressions,
			Clicks:      item.Clicks,
			Conversions: 0,
			Cost:        item.CostInLocalCurrency,
			Revenue:     0.0,
			Timestamp:   time.Now().UTC().String(),
		}
		processor.ProcessMetric(metric)
	}
}
