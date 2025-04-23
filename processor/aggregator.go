// processor/aggregator.go
package processor

import (
	"fmt"
	"time"

	"campaign-analytics/models"
	"campaign-analytics/storage"
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

	cpa := 0.0
	if m.Conversions > 0 {
		cpa = m.Cost / float64(m.Conversions)
	}

	fmt.Printf("Processed Campaign: %s | CTR: %.2f | ROAS: %.2f | CPA: %.2f\n", m.CampaignID, ctr, roas, cpa)

	// Retry insert up to 3 times on error (excluding dedup conflict)
	var err error
	for i := 0; i < 3; i++ {
		err = storage.InsertCampaignMetrics(m)
		if err == nil {
			break
		}
		fmt.Printf("Retrying DB insert (attempt %d) due to error: %v\n", i+1, err)
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		fmt.Printf("Final failure inserting into DB for %s: %v\n", m.CampaignID, err)
	}
}
