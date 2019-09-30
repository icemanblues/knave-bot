package karma

import (
	"github.com/gin-gonic/gin"
)

// BindRoutes bind handlers to router
func BindRoutes(karmaGroup *gin.RouterGroup, knaveGroup *gin.RouterGroup, karmaHandler Handler) {
	v1 := karmaGroup.Group("/v1")
	// team
	v1.GET("/team/:team", karmaHandler.TopKarma)
	// team user
	v1.GET("/team/:team/:user", karmaHandler.GetKarma)
	v1.PUT("/team/:team/:user", karmaHandler.AddKarma)
	v1.DELETE("/team/:team/:user", karmaHandler.DelKarma)

	// slack slash command integration
	slash := knaveGroup.Group("v1")
	slash.POST("/cmd/karma", karmaHandler.SlashKarma)
}
