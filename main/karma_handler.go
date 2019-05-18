package main

import (
	"github.com/gin-gonic/gin"
)

type KarmaHandler interface {
	GetKarma(c *gin.Context)
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

func NewKarmaHandler(kdb KarmaDB) *LiteKarmaHandler {
	return &LiteKarmaHandler{kdb}
}
