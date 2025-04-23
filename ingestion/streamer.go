package ingestion

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/processor"
)

// platforms is a mock set of ad platforms
var platforms = []string{"Meta", "Google", "LinkedIn", "TikTok"}

// SimulateStream generates random campaign metrics and sends them for processing
func SimulateStream() {
	go func() {
		for {
			metric := models.CampaignMetrics{
				CampaignID:  fmt.Sprintf("cmp-%d", rand.Intn(100)),
				Platform:    platforms[rand.Intn(len(platforms))],
				Impressions: rand.Intn(1000),
				Clicks:      rand.Intn(200),
				Conversions: rand.Intn(50),
				Cost:        float64(rand.Intn(5000)) / 100,
				Revenue:     float64(rand.Intn(10000)) / 100,
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			data, _ := json.Marshal(metric)
			fmt.Println("Ingested:", string(data))

			// Send the metric to aggregator for processing
			processor.ProcessMetric(metric)

			time.Sleep(2 * time.Second) // simulate delay
		}
	}()
}
