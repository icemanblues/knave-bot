package knave

import (
	"github.com/icemanblues/knave-bot/shakespeare"

	"github.com/gin-gonic/gin"
)

// Handler handler functions interface
type Handler interface {
	Insult(c *gin.Context)
	Compliment(c *gin.Context)
	SlashKnave(c *gin.Context)
}

// GinHandler an implementation using Gin
type GinHandler struct{}

// Insult handler function to generate an insult
func (g *GinHandler) Insult(c *gin.Context) {
	c.String(200, "%s", shakespeare.Insult())
}

// Compliment handler function to generate a complement
func (g *GinHandler) Compliment(c *gin.Context) {
	c.String(200, "%s", shakespeare.Compliment())
}

// SlashKnave handler function for slash-command `/knave`
func (g *GinHandler) SlashKnave(c *gin.Context) {
	c.JSON(200, gin.H{
		"text":          shakespeare.Insult(),
		"response_type": "in_channel",
	})
}

// NewHandler factory method
func NewHandler() *GinHandler {
	return &GinHandler{}
}
