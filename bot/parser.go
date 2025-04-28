package bot

import (
	"errors"
	"strings"
)

// Intent struct defines what user wants
type Intent struct {
	Metric    string
	Operation string
}

// Filters struct defines any filters user asked for
type Filters struct {
	Platform string
	FromDate string
	ToDate   string
}

// Very basic keyword-based parser
func ParsePrompt(prompt string) (Intent, Filters, error) {
	prompt = strings.ToLower(prompt)
	intent := Intent{}
	filters := Filters{}

	// Identify metric
	if strings.Contains(prompt, "roas") {
		intent.Metric = "roas"
	} else if strings.Contains(prompt, "ctr") {
		intent.Metric = "ctr"
	} else if strings.Contains(prompt, "spend") {
		intent.Metric = "spend"
	} else {
		return intent, filters, errors.New("unknown metric")
	}

	// Identify platform
	if strings.Contains(prompt, "google") {
		filters.Platform = "Google"
	} else if strings.Contains(prompt, "meta") {
		filters.Platform = "Meta"
	} else if strings.Contains(prompt, "linkedin") {
		filters.Platform = "LinkedIn"
	} else if strings.Contains(prompt, "tiktok") {
		filters.Platform = "TikTok"
	}

	// (optional) In real parsing, handle timeframes like "last week", "this quarter"
	filters.FromDate = "2024-01-01" // placeholder for now
	filters.ToDate = "2024-04-01"

	return intent, filters, nil
}
