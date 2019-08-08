package main

import (
	"fmt"

	"github.com/icemanblues/knave-bot/karma"
	"github.com/icemanblues/knave-bot/knave"
	"github.com/icemanblues/knave-bot/shakespeare"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// TODO: Move these to log lines
	fmt.Printf("Insult: %v\n", shakespeare.Insult())
	fmt.Printf("Compliment: %v\n", shakespeare.Compliment())

	// initialize database
	db, err := karma.InitDB()
	if err != nil {
		panic(err)
	}

	// create databinders
	kdb := karma.NewKdb(db)

	// create processors
	karmaProc := karma.NewProcessor(kdb)

	// create handlers
	knave := knave.NewHandler()
	karma := karma.NewHandler(karmaProc, kdb)

	// create gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	// map the routes to the handlers
	BindRoutes(r, knave, karma)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}
