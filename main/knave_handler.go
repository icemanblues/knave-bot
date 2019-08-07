package main

import "github.com/gin-gonic/gin"

// KnaveHandler handler functions interface
type KnaveHandler interface {
	Insult(c *gin.Context)
	Compliment(c *gin.Context)
	SlashKnave(c *gin.Context)
}

// GinKnaveHandler an implementation using Gin
type GinKnaveHandler struct{}

// Insult handler function to generate an insult
func (gkh *GinKnaveHandler) Insult(c *gin.Context) {
	c.String(200, "%s", Insult())
}

// Compliment handler function to generate a complement
func (gkh *GinKnaveHandler) Compliment(c *gin.Context) {
	c.String(200, "%s", Compliment())
}

// SlashKnave handler function for slash-command `/knave`
func (gkh *GinKnaveHandler) SlashKnave(c *gin.Context) {
	c.JSON(200, gin.H{
		"text":          Insult(),
		"response_type": "in_channel",
	})
}

// NewKnaveHandler factory method
func NewKnaveHandler() *GinKnaveHandler {
	return &GinKnaveHandler{}
}
