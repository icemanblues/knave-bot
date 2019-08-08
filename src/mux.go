package main

import (
	"github.com/icemanblues/knave-bot/karma"
	"github.com/icemanblues/knave-bot/knave"

	"github.com/gin-gonic/gin"
)

// BindRoutes bind handlers to router
func BindRoutes(r *gin.Engine, knave knave.Handler, karma karma.Handler) {
	// knave compliment and insult
	r.GET("/knavebot/v1/insult", knave.Insult)
	r.GET("/knavebot/v1/compliment", knave.Compliment)

	// karma
	r.GET("/karmabot/v1/:team/:user", karma.GetKarma)
	r.PUT("/karmabot/v1/:team/:user", karma.AddKarma)
	r.DELETE("/karmabot/v1/:team/:user", karma.DelKarma)

	// slack slash command integration
	r.POST("/knavebot/v1/cmd/knave", knave.SlashKnave)
	r.POST("/knavebot/v1/cmd/karma", karma.SlashKarma)
	// backwards compatibility knave
	r.POST("/knavebot/v1/insult", knave.SlashKnave)
}
