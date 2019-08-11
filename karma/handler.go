package karma

import (
	"strconv"

	"github.com/icemanblues/knave-bot/slack"

	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

// Handler interface for Karma handler
type Handler interface {
	GetKarma(c *gin.Context)
	AddKarma(c *gin.Context)
	DelKarma(c *gin.Context)
	SlashKarma(c *gin.Context)
}

// SQLiteHandler Karma Handler implementation using sqlite
type SQLiteHandler struct {
	kProc Processor
	kdb   DAO
}

// GetKarma handler method to read the current karma for an individual
func (lkh *SQLiteHandler) GetKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k, err := lkh.kdb.GetKarma(team, user)
	if err != nil {
		log.Error("Unable to look up karma.", team, user, err)
		c.String(500, "I'm sorry, we are not able to do that right now")
	}

	c.String(200, "%v", k)
}

// AddKarma handler method to add (or subtract) karma from an individual
func (lkh *SQLiteHandler) AddKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	d := c.Query("delta")
	delta, err := strconv.Atoi(d)
	if err != nil {
		log.Error("Not a valid integer.", team, user, delta, err)
		c.String(400, "Please pass a valid integer. %v", d)
		return
	}

	k, err := lkh.kdb.UpdateKarma(team, user, delta)
	if err != nil {
		log.Error("Unable to add or remove karma.", team, user, delta, err)
		c.String(500, "I'm sorry, we are not able to do that right now")
	}
	c.String(200, "%v", k)
}

// DelKarma handler method to delete (reset) karma to zer0
func (lkh *SQLiteHandler) DelKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k, err := lkh.kdb.DeleteKarma(team, user)
	if err != nil {
		log.Error("Unable to reset karma.", team, user, err)
		c.String(500, "I'm sorry, we are not able to do that right now")
	}

	c.String(200, "%v", k)
}

// SlashKarma handler method for the `/karma` slash-command
func (lkh *SQLiteHandler) SlashKarma(c *gin.Context) {
	data := &slack.CommandData{
		Command:      c.PostForm("command"),
		Text:         c.PostForm("text"),
		ResponseURL:  c.PostForm("response_url"),
		EnterpriseID: c.PostForm("enterprise_id"),
		TeamID:       c.PostForm("team_id"),
		ChannelID:    c.PostForm("channel_id"),
		UserID:       c.PostForm("user_id"),
	}

	response, err := lkh.kProc.Process(data)
	if err != nil {
		log.Error("Could not process a slack slash command.", data, err)
		response = slack.ErrorResponse("Oh no! Looks like we're experiencing some technical difficulties")
	}

	c.JSON(200, response)
}

// NewHandler factory method
func NewHandler(kProc Processor, kdb DAO) *SQLiteHandler {
	return &SQLiteHandler{
		kProc: kProc,
		kdb:   kdb,
	}
}
