package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println(Insult())
	fmt.Println(Compliment())

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
	r.GET("/knavebot/v1/insult", knave.Insult)
	r.GET("/knavebot/v1/compliment", knave.Compliment)
	// https://api.slack.com/slash-commands
	r.POST("/knavebot/v1/insult", knave.SlashKnave)

	r.GET("/karmabot/v1/:team/:user", karma.GetKarma)

	r.Run() // listen and serve on 0.0.0.0:8080
}
