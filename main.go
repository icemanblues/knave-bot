package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(Insult())
	fmt.Println(Compliment())

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/knavebot/v1/insult", func(c *gin.Context) {
		c.String(200, "%s", Insult())
	})

	r.GET("/knavebot/v1/compliment", func(c *gin.Context) {
		c.String(200, "%s", Compliment())
	})

	// https://api.slack.com/slash-commands
	r.POST("/knavebot/v1/insult", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"text":          Insult(),
			"response_type": "in_channel",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
