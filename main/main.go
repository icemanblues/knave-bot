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

	// knave compliment and insult
	r.GET("/knavebot/v1/insult", knave.Insult)
	r.GET("/knavebot/v1/compliment", knave.Compliment)

	// karma
	r.GET("/karmabot/v1/:team/:user", karma.GetKarma)
	r.PUT("/karmabot/v1/:team/:user", karma.AddKarma)
	r.DELETE("/karmabot/v1/:team/:user", karma.DelKarma)

	// slack slash command integration
	r.POST("/knavebot/v1/cmd/knave", knave.SlashKnave)
	r.POST("knavebot/v1/cmd/karma", karma.SlashKarma)
	// backwards compatibility knave
	r.POST("/knavebot/v1/insult", knave.SlashKnave)

	r.Run() // listen and serve on 0.0.0.0:8080
}
