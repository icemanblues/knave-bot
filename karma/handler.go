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
	TopKarma(c *gin.Context)
}

// SQLiteHandler Karma Handler implementation using sqlite
type SQLiteHandler struct {
	proc     Processor
	dao      DAO
	dailyDao DailyDao
}

// GetKarma handler method to read the current karma for an individual
func (h *SQLiteHandler) GetKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k, err := h.dao.GetKarma(team, user)
	if err != nil {
		log.Errorf("Unable to lookup karma. %v %v %v", team, user, err)
		c.String(500, err.Error())
		return
	}

	c.String(200, "%v", k)
}

// AddKarma handler method to add (or subtract) karma from an individual
func (h *SQLiteHandler) AddKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	d := c.Query("delta")
	delta, err := strconv.Atoi(d)
	if err != nil {
		log.Errorf("Not a valid integer. %v %v %v %v", team, user, delta, err)
		c.String(400, "Please pass a valid integer. %v", d)
		return
	}

	k, err := h.dao.UpdateKarma(team, user, delta)
	if err != nil {
		log.Errorf("Unable to add or remove karma. %v %v %v %v", team, user, delta, err)
		c.String(500, err.Error())
		return
	}
	c.String(200, "%v", k)
}

// DelKarma handler method to delete (reset) karma to zer0
func (h *SQLiteHandler) DelKarma(c *gin.Context) {
	team := c.Param("team")
	user := c.Param("user")

	k, err := h.dao.DeleteKarma(team, user)
	if err != nil {
		log.Errorf("Unable to reset karma. %v %v %v", team, user, err)
		c.String(500, err.Error())
		return
	}

	c.String(200, "%v", k)
}

// TopKarma returns the top n users for a given team
func (h *SQLiteHandler) TopKarma(c *gin.Context) {
	team := c.Param("team")

	top := c.Query("top")
	n, err := strconv.Atoi(top)
	if err != nil {
		c.String(400, "Please pass a valid integer. %v", n)
		return
	}
	if n <= 0 {
		c.String(400, "Please pass a positive non-zero integer. %v", n)
		return
	}

	topUsers, err := h.dao.Top(team, n)
	c.JSON(200, topUsers)
}

var responseUnknownError = slack.ErrorResponse("Oh no! Looks like we're experiencing some technical difficulties")

// SlashKarma handler method for the `/karma` slash-command
func (h *SQLiteHandler) SlashKarma(c *gin.Context) {
	data := &slack.CommandData{
		Command:      c.PostForm("command"),
		Text:         c.PostForm("text"),
		ResponseURL:  c.PostForm("response_url"),
		EnterpriseID: c.PostForm("enterprise_id"),
		TeamID:       c.PostForm("team_id"),
		ChannelID:    c.PostForm("channel_id"),
		UserID:       c.PostForm("user_id"),
	}

	response, err := h.proc.Process(data)
	if err != nil {
		log.Errorf("Could not process a slack slash command. %v %v", data, err)
		response = responseUnknownError
	}

	c.JSON(200, response)

	// async log the usage
	go func() {
		if err := h.dao.Usage(data, response); err != nil {
			log.Errorf("Unable to log karma usage %v", err)
		}
	}()
}

// NewHandler factory method
func NewHandler(proc Processor, dao DAO) *SQLiteHandler {
	return &SQLiteHandler{
		proc: proc,
		dao:  dao,
	}
}
