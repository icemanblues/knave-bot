package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(Generate())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/insult", func(c *gin.Context) {
		c.String(200, "%s", Generate())
	})

	// https://api.slack.com/slash-commands
	r.POST("/insult", func(c *gin.Context) {
		// read body as JSON
		// check the text for the user id

		c.JSON(200, gin.H{
			"text":          Generate(),
			"response_type": "in_channel",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
