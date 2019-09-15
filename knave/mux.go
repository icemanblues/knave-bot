package knave

import (
	"github.com/gin-gonic/gin"
)

// BindRoutes bind handlers to router
func BindRoutes(r *gin.RouterGroup, knave Handler) {
	v1 := r.Group("/v1")
	v1.GET("/insult", knave.Insult)
	v1.GET("/compliment", knave.Compliment)

	// slack slash command integration
	v1.POST("/cmd/knave", knave.SlashKnave)
}
