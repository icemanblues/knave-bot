package knave

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icemanblues/knave-bot/slack"
	"github.com/stretchr/testify/assert"
)

func setup(h *GinHandler) *gin.Engine {
	r := gin.Default()

	r.GET("/insult", h.Insult)
	r.GET("/compliment", h.Compliment)
	r.POST("/slash/knave", h.SlashKnave)

	return r
}

func TestSlashKnave(t *testing.T) {
	h := NewHandler()
	r := setup(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/slash/knave", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var body slack.Response
	json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, slack.ResponseType.InChannel, body.ResponseType)
	assert.True(t, len(body.Text) > 4)
	assert.True(t, strings.HasPrefix(body.Text, "Thou"))
}
