package main

import "github.com/gin-gonic/gin"

// BindRoutes bind handlers to router
func BindRoutes(r *gin.Engine, knave KnaveHandler, karma KarmaHandler) {
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
