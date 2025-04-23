// processor/aggregator.go
package processor

import (
	"campaign-analytics/models"
	"campaign-analytics/storage"
	"fmt"
)

// ProcessMetric calculates derived metrics and stores them in DB
func ProcessMetric(m models.CampaignMetrics) {
	// Calculate derived metrics
	ctr := 0.0
	if m.Impressions > 0 {
		ctr = float64(m.Clicks) / float64(m.Impressions)
	}

	roas := 0.0
	if m.Cost > 0 {
		roas = m.Revenue / m.Cost
	}

	// Log derived metrics
	fmt.Printf("Processed Campaign: %s | CTR: %.2f | ROAS: %.2f\n", m.CampaignID, ctr, roas)

	// Persist raw metric data into the database
	err := storage.InsertCampaignMetrics(m)
	if err != nil {
		fmt.Printf("Failed to insert into DB: %v\n", err)
	}
}
