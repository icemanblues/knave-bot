package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Printf("Insult: %v\n", Insult())
	fmt.Printf("Compliment: %v\n", Compliment())

	// initialize database
	db, err := InitDB()
	if err != nil {
		panic(err)
	}

	// create databinders
	kdb := NewKdb(db)

	// create handlers
	knave := NewKnaveHandler()
	karma := NewKarmaHandler(kdb)

	// create gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	// bind handlers to router
	// https://api.slack.com/slash-commands

	r.GET("/knavebot/v1/insult", knave.Insult)
	r.GET("/knavebot/v1/compliment", knave.Compliment)

	// knave
	r.POST("/knavebot/v1/insult", knave.SlashKnave)

	// karma
	r.GET("/karmabot/v1/:team/:user", karma.GetKarma)
	r.PUT("/karmabot/v1/:team/:user", karma.AddKarma)
	r.DELETE("/karmabot/v1/:team/:user", karma.DelKarma)

	r.Run() // listen and serve on 0.0.0.0:8080
}
