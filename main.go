// main.go
package main

import (
	"log"

	"campaign-analytics/api"
	"campaign-analytics/ingestion"
	"campaign-analytics/storage"
)

// main initializes all services: DB, Redis, ingestion, and starts the API server
func main() {
	// Step 1: Initialize database connection
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Step 2: Initialize Redis cache
	if err := storage.InitRedis(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Step 3: Start simulating data ingestion
	ingestion.SimulateStream()

	// Step 4: Start REST API server
	r := api.InitRouter()
	r.Run(":8080")
}
