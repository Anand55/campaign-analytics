package bot

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// SearchNearestCampaign finds the campaign closest to prompt embedding
func SearchNearestCampaign(embedding []float32) (string, error) {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return "", err
	}
	defer db.Close()

	vectorString := fmt.Sprintf("[%v]", embedding)

	var campaignID string
	query := `
		SELECT campaign_id
		FROM campaign_embeddings
		ORDER BY embedding <-> $1
		LIMIT 1;
	`
	err = db.QueryRow(query, vectorString).Scan(&campaignID)
	if err != nil {
		return "", err
	}

	return campaignID, nil
}
