package main

import "github.com/gin-gonic/gin"

type KnaveHandler interface {
	Insult(c *gin.Context)
	Compliment(c *gin.Context)
	SlashKnave(c *gin.Context)
}

type GinKnaveHandler struct{}

func (gkh *GinKnaveHandler) Insult(c *gin.Context) {
	c.String(200, "%s", Insult())
}

func (gkh *GinKnaveHandler) Compliment(c *gin.Context) {
	c.String(200, "%s", Compliment())
}

func (gkh *GinKnaveHandler) SlashKnave(c *gin.Context) {
	c.JSON(200, gin.H{
		"text":          Insult(),
		"response_type": "in_channel",
	})
}

func NewKnaveHandler() *GinKnaveHandler {
	return &GinKnaveHandler{}
}
