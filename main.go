package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(Generate())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

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

	// OLD GET AND POST
	// These are only here until the cutover is complete
	r.GET("/api/v1/knave/insult", func(c *gin.Context) {
		c.String(200, "%s", Generate())
	})
	r.POST("/api/v1/knave/insult", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"text":          Generate(),
			"response_type": "in_channel",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
