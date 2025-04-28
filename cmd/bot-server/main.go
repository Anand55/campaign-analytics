package main

import (
	"campaign-analytics/bot"
	"fmt"
)

func main() {
	fmt.Println("✅ Bot Server started on port 8081")
	bot.StartBotServer()
}
