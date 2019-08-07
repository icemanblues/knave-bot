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

	// create processors
	karmaProc := NewKdbProcessor(kdb)

	// create handlers
	knave := NewKnaveHandler()
	karma := NewKarmaHandler(karmaProc, kdb)

	// create gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	// map the routes to the handlers
	BindRoutes(r, knave, karma)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}
