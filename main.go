package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(Generate())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/knavebot/v1/insult", func(c *gin.Context) {
		c.String(200, "%s", Generate())
	})

	// https://api.slack.com/slash-commands
	r.POST("/knavebot/v1/insult", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"text":          Generate(),
			"response_type": "in_channel",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
