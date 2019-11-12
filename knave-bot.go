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

func logger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetReportCaller(true) // This could have performance impact
}

// InitKarma initializes the components and wires them together, for Karma and Knave bot
func initKarma(insult, compliment shakespeare.Generator, config karma.ProcConfig, dao karma.DAO, dailyDao karma.DailyDao) (knave.Handler, karma.Handler) {
	karmaProc := karma.NewProcessor(config, dao, dailyDao, insult, compliment)

	knave := knave.NewHandler(insult, compliment)
	karma := karma.NewHandler(karmaProc, dao)

	return knave, karma
}

func initGin() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	return r
}

func main() {
	// initialize logger
	logger()
	log.Infof("Insult    : %v", shakespeare.Insult())
	log.Infof("Compliment: %v", shakespeare.Compliment())

	// initialize database
	db, err := karma.InitDB("/var/lib/sqlite/karma.db")
	if err != nil {
		log.Panic("Unable to initialize the database", err)
		panic(err)
	}
	dao := karma.NewDao(db)
	dailydao := karma.NewDailyDao(db)

	procConfig := karma.DefaultConfig
	knaveHandler, karmaHandler := initKarma(shakespeare.InsultGenerator, shakespeare.ComplimentGenerator, procConfig, dao, dailydao)

	r := initGin()
	BindRoutes(r, knaveHandler, karmaHandler)

	// listen and serve on 0.0.0.0:8080
	r.Run()
}
