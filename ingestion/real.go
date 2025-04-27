// // ingestion/real.go
package ingestion

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"

// 	"campaign-analytics/models"
// 	"campaign-analytics/processor"
// )

// // StartRealFetcher pulls real campaign data from Meta Ads API every minute
// func StartRealFetcher() {
// 	fmt.Println("[INFO] Using real data fetcher: Meta Ads API")
// 	ticker := time.NewTicker(60 * time.Second)
// 	defer ticker.Stop()

// 	for {
// 		fetchMetaCampaignData()
// 		<-ticker.C
// 	}
// }

// // fetchMetaCampaignData makes a GET call to Meta Ads Insights API
// func fetchMetaCampaignData() {
// 	token := os.Getenv("META_ACCESS_TOKEN")
// 	adAccountID := os.Getenv("META_AD_ACCOUNT_ID")

// 	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/insights?fields=campaign_name,impressions,clicks,spend,date_stop&access_token=%s", adAccountID, token)

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Printf("[ERROR] Failed to fetch Meta data: %v\n", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		fmt.Printf("[WARN] Meta API returned %d\n", resp.StatusCode)
// 		return
// 	}

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	var response struct {
// 		Data []map[string]string `json:"data"`
// 	}
// 	if err := json.Unmarshal(body, &response); err != nil {
// 		fmt.Printf("[ERROR] Failed to parse Meta response: %v\n", err)
// 		return
// 	}

// 	for _, item := range response.Data {
// 		impressions, _ := strconv.Atoi(item["impressions"])
// 		clicks, _ := strconv.Atoi(item["clicks"])
// 		cost, _ := strconv.ParseFloat(item["spend"], 64)
// 		timestamp, _ := time.Parse("2006-01-02", item["date_stop"])

// 		metric := models.CampaignMetrics{
// 			CampaignID:  item["campaign_name"],
// 			Platform:    "Meta",
// 			Impressions: impressions,
// 			Clicks:      clicks,
// 			Conversions: 0, // not available
// 			Cost:        cost,
// 			Revenue:     0.0, // not available
// 			Timestamp:   timestamp.UTC().String(),
// 		}
// 		processor.ProcessMetric(metric)
// 	}
// }
