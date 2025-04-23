// storage/db.go
package storage

import (
	"campaign-analytics/models"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// DB is a shared global database connection object
var DB *sql.DB

// InitDB initializes the PostgreSQL connection
func InitDB() error {
	connStr := "host=postgres port=5432 user=postgres dbname=campaigns password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	DB = db
	fmt.Println("Connected to Postgres")
	return nil
}

// InsertCampaignMetrics inserts a metrics record into the DB
func InsertCampaignMetrics(m models.CampaignMetrics) error {
	query := `INSERT INTO campaign_metrics
		(campaign_id, platform, impressions, clicks, conversions, cost, revenue, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := DB.Exec(query,
		m.CampaignID,
		m.Platform,
		m.Impressions,
		m.Clicks,
		m.Conversions,
		m.Cost,
		m.Revenue,
		m.Timestamp,
	)
	return err
}
