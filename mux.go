package main

import (
	"github.com/icemanblues/knave-bot/karma"
	"github.com/icemanblues/knave-bot/knave"

	"github.com/gin-gonic/gin"
)

// BindRoutes bind handlers to router
func BindRoutes(r *gin.Engine, knaveHandler knave.Handler, karmaHandler karma.Handler) {
	knaveRouter := r.Group("/knavebot")
	knave.BindRoutes(knaveRouter, knaveHandler)

	karmaRouter := r.Group("/karmabot")
	karma.BindRoutes(karmaRouter, knaveRouter, karmaHandler)
}
