package bot

import "github.com/gin-gonic/gin"

func StartBotServer() {
	r := gin.Default()
	r.POST("/prompt", PromptHandler)
	r.Run(":8081")
}
