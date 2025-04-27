// ingestion/meta.go
package ingestion

import (
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

// FetchMetaCampaigns pulls campaign insights from Meta Ads API
func FetchMetaCampaigns() {
	fmt.Println("[META] Fetching campaign data from Meta Ads API...")

	token := os.Getenv("META_ACCESS_TOKEN")
	adAccountID := os.Getenv("META_AD_ACCOUNT_ID")

	if token == "" || adAccountID == "" {
		fmt.Println("[META] Missing Meta API credentials.")
		return
	}

	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/insights?fields=campaign_name,impressions,clicks,spend,date_stop&level=campaign&date_preset=last_30d&access_token=%s", adAccountID, token)

	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[META] Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[META] API returned non-200: %d\n", resp.StatusCode)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var response struct {
		Data []struct {
			CampaignName string `json:"campaign_name"`
			Impressions  string `json:"impressions"`
			Clicks       string `json:"clicks"`
			Spend        string `json:"spend"`
			DateStop     string `json:"date_stop"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Printf("[META] Failed to parse response: %v\n", err)
		return
	}

	for _, item := range response.Data {
		impressions := atoi(item.Impressions)
		clicks := atoi(item.Clicks)
		spend, _ := strconv.ParseFloat(item.Spend, 64)
		timestamp, _ := time.Parse("2006-01-02", item.DateStop)

		metric := models.CampaignMetrics{
			CampaignID:  fmt.Sprintf("m-%s", item.CampaignName),
			Platform:    "Meta",
			Impressions: impressions,
			Clicks:      clicks,
			Conversions: 0,
			Cost:        spend,
			Revenue:     0.0,
			Timestamp:   timestamp.UTC().String(),
		}
		processor.ProcessMetric(metric)
	}
}
