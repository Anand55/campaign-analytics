// ingestion/google.go
package ingestion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/processor"
)

// FetchGoogleCampaigns pulls campaign insights from Google Ads API via REST
func FetchGoogleCampaigns() {
	fmt.Println("[GOOGLE] Fetching campaign data from Google Ads API (REST)...")

	accessToken := os.Getenv("GOOGLE_ADS_ACCESS_TOKEN")
	customerID := os.Getenv("GOOGLE_ADS_CUSTOMER_ID")

	if accessToken == "" || customerID == "" {
		fmt.Println("[GOOGLE] Missing Google Ads API credentials.")
		return
	}

	url := fmt.Sprintf("https://googleads.googleapis.com/v16/customers/%s/googleAds:search", customerID)

	query := map[string]interface{}{
		"query": `SELECT campaign.id, campaign.name, metrics.impressions, metrics.clicks, metrics.cost_micros FROM campaign WHERE campaign.status = 'ENABLED' LIMIT 10`,
	}
	payload, _ := json.Marshal(query)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[GOOGLE] Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[GOOGLE] API returned non-200: %d\n", resp.StatusCode)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var response struct {
		Results []struct {
			Campaign struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"campaign"`
			Metrics struct {
				Impressions string `json:"impressions"`
				Clicks      string `json:"clicks"`
				CostMicros  string `json:"costMicros"`
			} `json:"metrics"`
		} `json:"results"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Printf("[GOOGLE] Failed to parse response: %v\n", err)
		return
	}

	for _, row := range response.Results {
		impressions := atoi(row.Metrics.Impressions)
		clicks := atoi(row.Metrics.Clicks)
		costMicros := atoi(row.Metrics.CostMicros)
		cost := float64(costMicros) / 1_000_000

		metric := models.CampaignMetrics{
			CampaignID:  fmt.Sprintf("g-%s", row.Campaign.Id),
			Platform:    "Google",
			Impressions: impressions,
			Clicks:      clicks,
			Conversions: 0,
			Cost:        cost,
			Revenue:     0.0,
			Timestamp:   time.Now().UTC().String(),
		}
		processor.ProcessMetric(metric)
	}
}

// atoi safely parses string to int
func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}
