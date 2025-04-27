// ingestion/dispatcher.go
package ingestion

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// StartRealFetcher determines which APIs to call based on env flags
func StartRealFetcher() {
	fmt.Println("[DISPATCHER] Starting real API ingestion mode...")

	enabledSources := strings.Split(os.Getenv("ENABLED_SOURCES"), ",")
	sourceMap := make(map[string]bool)
	for _, src := range enabledSources {
		sourceMap[strings.TrimSpace(strings.ToLower(src))] = true
	}

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	callAll := func() {
		if sourceMap["meta"] {
			FetchMetaCampaigns()
		}
		if sourceMap["google"] {
			FetchGoogleCampaigns()
		}
		if sourceMap["tiktok"] {
			FetchTiktokCampaigns()
		}
		if sourceMap["linkedin"] {
			FetchLinkedInCampaigns()
		}
	}

	// Initial call
	callAll()

	for {
		<-ticker.C
		callAll()
	}
}
