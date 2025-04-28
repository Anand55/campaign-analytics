// main.go
package main

import (
	"fmt"
	"os"
	"time"

	"campaign-analytics/api"
	"campaign-analytics/ingestion"
	"campaign-analytics/storage"
)

func main() {
	// Initialize Postgres
	if err := storage.InitDB(); err != nil {
		fmt.Println("[ERROR] Failed to connect to DB:", err)
		os.Exit(1)
	}

	// Initialize Redis
	if err := storage.InitRedis(); err != nil {
		fmt.Println("[ERROR] Failed to connect to Redis:", err)
		os.Exit(1)
	}
	fmt.Println("[INFO] Connected to Redis")

	// Decide ingestion mode
	mode := os.Getenv("DATA_SOURCE")
	if mode == "real" {
		fmt.Println("[BOOT] Running in REAL ingestion mode (Meta, Google, TikTok, LinkedIn)")
		go ingestion.StartRealFetcher()
	} else {
		fmt.Println("[BOOT] Running in FAKE data simulation mode")
		go ingestion.StartSimulator()
	}

	// Small delay to ensure ingestion is warmed up
	time.Sleep(1 * time.Second)

	// Start API server
	r := api.InitRouter()
	fmt.Println("[INFO] API Server started at http://localhost:8080")
	r.Run(":8080")
}
