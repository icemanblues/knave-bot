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
type GinHandler struct {
	insult     shakespeare.Generator
	compliment shakespeare.Generator
}

// Insult handler function to generate an insult
func (g *GinHandler) Insult(c *gin.Context) {
	c.String(200, "%s", g.insult.Sentence())
}

// Compliment handler function to generate a complement
func (g *GinHandler) Compliment(c *gin.Context) {
	c.String(200, "%s", g.compliment.Sentence())
}

// SlashKnave handler function for slash-command `/knave`
func (g *GinHandler) SlashKnave(c *gin.Context) {
	c.JSON(200, gin.H{
		"text":          g.insult.Sentence(),
		"response_type": "in_channel",
	})
}

// NewHandler factory method
func NewHandler(insult, compliment shakespeare.Generator) *GinHandler {
	return &GinHandler{
		insult:     insult,
		compliment: compliment,
	}
}
