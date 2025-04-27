// ingestion/tiktok.go
package ingestion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/processor"
)

// FetchTiktokCampaigns pulls campaign insights from TikTok Ads API
func FetchTiktokCampaigns() {
	fmt.Println("[TIKTOK] Fetching campaign data from TikTok Marketing API...")

	token := os.Getenv("TIKTOK_ACCESS_TOKEN")
	advertiserID := os.Getenv("TIKTOK_ADVERTISER_ID")

	if token == "" || advertiserID == "" {
		fmt.Println("[TIKTOK] Missing TikTok API credentials.")
		return
	}

	url := "https://business-api.tiktok.com/open_api/v1.3/report/integrated/get/"

	payload := map[string]interface{}{
		"advertiser_id": advertiserID,
		"report_type":   "BASIC",
		"dimensions":    []string{"campaign_id"},
		"metrics":       []string{"impressions", "clicks", "spend"},
		"data_level":    "CAMPAIGN",
		"start_date":    time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
		"end_date":      time.Now().Format("2006-01-02"),
	}

	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Access-Token", token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[TIKTOK] Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[TIKTOK] API returned non-200: %d\n", resp.StatusCode)
		return
	}

	respBody, _ := ioutil.ReadAll(resp.Body)
	var response struct {
		Data struct {
			List []struct {
				CampaignID  string  `json:"campaign_id"`
				Impressions int     `json:"impressions"`
				Clicks      int     `json:"clicks"`
				Spend       float64 `json:"spend"`
			} `json:"list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		fmt.Printf("[TIKTOK] Failed to parse response: %v\n", err)
		return
	}

	for _, item := range response.Data.List {
		metric := models.CampaignMetrics{
			CampaignID:  fmt.Sprintf("t-%s", item.CampaignID),
			Platform:    "TikTok",
			Impressions: item.Impressions,
			Clicks:      item.Clicks,
			Conversions: 0,
			Cost:        item.Spend,
			Revenue:     0.0,
			Timestamp:   time.Now().UTC().String(),
		}
		processor.ProcessMetric(metric)
	}
}
