package models

// CampaignMetrics defines the structure for ad campaign analytics data.
type CampaignMetrics struct {
	CampaignID  string  `json:"campaign_id" db:"campaign_id"`
	Platform    string  `json:"platform" db:"platform"`
	Impressions int     `json:"impressions" db:"impressions"`
	Clicks      int     `json:"clicks" db:"clicks"`
	Conversions int     `json:"conversions" db:"conversions"`
	Cost        float64 `json:"cost" db:"cost"`
	Revenue     float64 `json:"revenue" db:"revenue"`
	Timestamp   string  `json:"timestamp" db:"timestamp"`
}
