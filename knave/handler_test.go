package knave

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icemanblues/knave-bot/shakespeare"
	"github.com/icemanblues/knave-bot/slack"
	"github.com/stretchr/testify/assert"
)

func setupHandler() GinHandler {
	return NewHandler(shakespeare.New("insult", "", nil),
		shakespeare.New("compliment", "", nil))
}

func setupGin(h GinHandler) *gin.Engine {
	r := gin.Default()
	g := r.Group("/knavebot")
	BindRoutes(g, h)

	return r
}

func TestInsult(t *testing.T) {
	h := setupHandler()
	r := setupGin(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/knavebot/v1/insult", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "insult", w.Body.String())
}

func TestCompliment(t *testing.T) {
	h := setupHandler()
	r := setupGin(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/knavebot/v1/compliment", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "compliment", w.Body.String())
}

func TestSlashKnave(t *testing.T) {
	h := setupHandler()
	r := setupGin(h)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/knavebot/v1/cmd/knave", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var body slack.Response
	json.Unmarshal(w.Body.Bytes(), &body)

	assert.Equal(t, slack.ResponseType.InChannel, body.ResponseType)
	assert.True(t, len(body.Text) > 4)
	assert.Equal(t, "insult", body.Text)
}
