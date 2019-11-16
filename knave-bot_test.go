package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/icemanblues/knave-bot/karma"
	"github.com/icemanblues/knave-bot/shakespeare"
	"github.com/stretchr/testify/assert"
)

const testDB = "var/test/functional.db"

func setupDB(datasource string) (karma.DAO, error) {
	if err := os.RemoveAll(datasource); err != nil {
		return nil, err
	}

	db, err := karma.InitDB(datasource)
	if err != nil {
		return nil, err
	}

	return karma.NewDao(db), nil
}

func setup(t *testing.T) *gin.Engine {
	dao, err := setupDB(testDB)
	assert.Nil(t, err)

	insult := shakespeare.New("insult", "", nil)
	compliment := shakespeare.New("compliment", "", nil)
	knave, karma := initKarma(insult, compliment, karma.DefaultConfig, dao)
	r := initGin()
	BindRoutes(r, knave, karma)
	return r
}

func TestKnaveInsult(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/knavebot/v1/insult", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "insult", w.Body.String())
}

func TestKnaveCompliment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/knavebot/v1/compliment", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "compliment", w.Body.String())
}

func TestKnaveSlashCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/knavebot/v1/cmd/knave", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// it will return the insult wrapped in a slack response
	// assert.Equal(t, "insult", w.Body.String())
}

func TestKarmaPut(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/karmabot/v1/team/team1/playerB?delta=3", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "3", w.Body.String())
}

func TestKarmaDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/karmabot/v1/team/team1/playerD", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "0", w.Body.String())
}

func TestKarmaGet(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/karmabot/v1/team/team1/playerA", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "0", w.Body.String())
}

func TestKarmaSlashCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping functional test")
	}

	r := setup(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/knavebot/v1/cmd/karma", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// TODO: check that it is returning the help message
}
