package main

import (
	"os"

	"github.com/icemanblues/knave-bot/karma"
	"github.com/icemanblues/knave-bot/knave"
	"github.com/icemanblues/knave-bot/shakespeare"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	// initialize logger
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true) // This could have performance impact

	log.Infof("Insult: %v", shakespeare.Insult())
	log.Infof("Compliment: %v", shakespeare.Compliment())

	// initialize database
	db, err := karma.InitDB()
	if err != nil {
		log.Panic("Unable to initialize the database", err)
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
