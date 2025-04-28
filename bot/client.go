package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// QueryAnalyticsBackend fetches analytics for a given campaign ID
func QueryAnalyticsBackend(campaignID string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	// Form API URL
	url := fmt.Sprintf("http://campaign-analytics-app:8080/campaign/%s/insights", campaignID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
