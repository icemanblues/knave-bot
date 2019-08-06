package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KarmaHandler interface {
	GetKarma(c *gin.Context)
	AddKarma(c *gin.Context)
	DelKarma(c *gin.Context)
	SlashKarma(c *gin.Context)
}

type LiteKarmaHandler struct {
	kdb KarmaDB
}

func (lkh *LiteKarmaHandler) GetKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k := lkh.kdb.GetKarma(team, user)

	c.String(200, "%v", k)
}

func (lkh *LiteKarmaHandler) AddKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	d := c.Query("delta")
	delta, err := strconv.Atoi(d)
	if err != nil {
		c.String(400, "Please pass a valid integer. %v", d)
	}

	k := lkh.kdb.UpdateKarma(team, user, delta)

	c.String(200, "%v", k)
}

func (lkh *LiteKarmaHandler) DelKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k := lkh.kdb.DeleteKarma(team, user)

	c.String(200, "%v", k)
}

func (lkh *LiteKarmaHandler) SlashKarma(c *gin.Context) {
	data := &CommandData{
		Command:      c.PostForm("command"),
		Text:         c.PostForm("text"),
		ResponseURL:  c.PostForm("response_url"),
		EnterpriseID: c.PostForm("enterprise_id"),
		TeamID:       c.PostForm("team_id"),
		ChannelID:    c.PostForm("channel_id"),
		UserID:       c.PostForm("user_id"),
	}
	fmt.Printf("Command Data: %+v\n", data)

	response := ErrorResponse("Hello, this is a work-in-progress. Ask roland for more details")
	c.JSON(200, response)
}

func NewKarmaHandler(kdb KarmaDB) *LiteKarmaHandler {
	return &LiteKarmaHandler{kdb}
}
