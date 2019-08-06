package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type KarmaHandler interface {
	GetKarma(c *gin.Context)
	AddKarma(c *gin.Context)
	DelKarma(c *gin.Context)
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

func NewKarmaHandler(kdb KarmaDB) *LiteKarmaHandler {
	return &LiteKarmaHandler{kdb}
}
