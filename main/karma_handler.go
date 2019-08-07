package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// KarmaHandler interface for Karma handler
type KarmaHandler interface {
	GetKarma(c *gin.Context)
	AddKarma(c *gin.Context)
	DelKarma(c *gin.Context)
	SlashKarma(c *gin.Context)
}

// LiteKarmaHandler Karma Handler implementation using sqlite
type LiteKarmaHandler struct {
	kProc KarmaProcessor
	kdb   KarmaDB
}

// GetKarma handler method to read the current karma for an individual
func (lkh *LiteKarmaHandler) GetKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k := lkh.kdb.GetKarma(team, user)

	c.String(200, "%v", k)
}

// AddKarma handler method to add (or subtract) karma from an individual
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

// DelKarma handler method to delete (reset) karma to zer0
func (lkh *LiteKarmaHandler) DelKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k := lkh.kdb.DeleteKarma(team, user)

	c.String(200, "%v", k)
}

// SlashKarma handler method for the `/karma` slash-command
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

	response, err := lkh.kProc.Process(data)
	if err != nil {
		response = ErrorResponse("Oh no! Looks like we're experiencing some technical difficulties")
	}
	fmt.Printf("Response: %v\n", response.Text)

	c.JSON(200, response)
}

// NewKarmaHandler factory method
func NewKarmaHandler(kProc KarmaProcessor, kdb KarmaDB) *LiteKarmaHandler {
	return &LiteKarmaHandler{
		kProc: kProc,
		kdb:   kdb,
	}
}
