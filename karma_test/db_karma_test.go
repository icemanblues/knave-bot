package karma_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/icemanblues/knave-bot/karma"
	"github.com/icemanblues/knave-bot/slack"
	"github.com/stretchr/testify/assert"
)

const testDB = "var/test/test.db"

func setupDB(datasource string) (*sql.DB, karma.DAO, error) {
	if err := os.RemoveAll(datasource); err != nil {
		return nil, nil, err
	}

	db, err := karma.InitDB(datasource)
	if err != nil {
		return nil, nil, err
	}

	return db, karma.NewDao(db), nil
}

func rowCountKarma(t *testing.T, db *sql.DB) int {
	row := db.QueryRow("SELECT count(*) FROM karma")
	var rowCount int
	err := row.Scan(&rowCount)
	assert.Nil(t, err)
	return rowCount
}

func rowCountUsage(t *testing.T, db *sql.DB) int {
	row := db.QueryRow("SELECT count(*) FROM usage")
	var rowCount int
	err := row.Scan(&rowCount)
	assert.Nil(t, err)
	return rowCount
}

func TestGetKarma(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	rowCount := rowCountKarma(t, db)
	assert.Zero(t, rowCount)

	k, err := dao.GetKarma("nycfc", "ring")
	assert.Nil(t, err)
	assert.Zero(t, k)
}

func TestUpdateKarma(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	// confirm zeror rows
	rc := rowCountKarma(t, db)
	assert.Zero(t, rc)

	// confirm this user has 0 karma
	k, err := dao.GetKarma("nycfc", "ring")
	assert.Nil(t, err)
	assert.Zero(t, k)

	// add karma, this should create a row (INSERT)
	k, err = dao.UpdateKarma("nycfc", "ring", 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, k)

	rc = rowCountKarma(t, db)
	assert.Equal(t, 1, rc)

	// confirm that the user has 2 karma
	k, err = dao.GetKarma("nycfc", "ring")
	assert.Nil(t, err)
	assert.Equal(t, 2, k)

	// add karma, this should update the existing row (UPDATE)
	k, err = dao.UpdateKarma("nycfc", "ring", 3)
	assert.Nil(t, err)
	assert.Equal(t, 5, k)

	rc = rowCountKarma(t, db)
	assert.Equal(t, 1, rc)

	// confirm that the user has 5 karma
	k, err = dao.GetKarma("nycfc", "ring")
	assert.Nil(t, err)
	assert.Equal(t, 5, k)
}

func TestDeleteKarma(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	rc := rowCountKarma(t, db)
	assert.Zero(t, rc)

	k, err := dao.UpdateKarma("nycfc", "maxi", 10)
	assert.Nil(t, err)
	assert.Equal(t, 10, k)

	k, err = dao.UpdateKarma("nycfc", "heber", 9)
	assert.Nil(t, err)
	assert.Equal(t, 9, k)

	rc = rowCountKarma(t, db)
	assert.Equal(t, 2, rc)

	k, err = dao.DeleteKarma("nycfc", "maxi")
	assert.Nil(t, err)
	assert.Equal(t, 0, k)

	rc = rowCountKarma(t, db)
	assert.Equal(t, 1, rc)

	k, err = dao.GetKarma("nycfc", "maxi")
	assert.Nil(t, err)
	assert.Zero(t, k)
}

// TODO: This only tests the inserts. doesn't confirm what is written
func TestUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	rowCount := rowCountUsage(t, db)
	assert.Zero(t, rowCount)

	// insert default values
	cd := slack.CommandData{}
	res := slack.Response{}
	dao.Usage(cd, res)
	rowCount = rowCountUsage(t, db)
	assert.Equal(t, 1, rowCount)

	// insert default values again (dupe) and it should work
	dao.Usage(cd, res)
	rowCount = rowCountUsage(t, db)
	assert.Equal(t, 2, rowCount)

	// insert direct message
	cd = slack.CommandData{
		Command:      "karma",
		Text:         "me",
		EnterpriseID: "enterprise",
		TeamID:       "team",
		ChannelID:    "channel",
		UserID:       "user",
	}
	res = slack.DirectResponse("slack.DirectResponse", "")
	dao.Usage(cd, res)

	rowCount = rowCountUsage(t, db)
	assert.Equal(t, 3, rowCount)

	// insert no attachments
	cd.Text = "status"
	res = slack.ChannelResponse("slack.ChannelResponse")
	dao.Usage(cd, res)

	rowCount = rowCountUsage(t, db)
	assert.Equal(t, 4, rowCount)

	// insert with attachments
	cd.Text = "attachments"
	res = slack.ChannelAttachmentsResponse(
		"slack.ChannelAttachmentsResponse",
		"attachments")
	dao.Usage(cd, res)

	rowCount = rowCountUsage(t, db)
	assert.Equal(t, 5, rowCount)
}

func TestTopKarma(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	rowCount := rowCountKarma(t, db)
	assert.Zero(t, rowCount)

	// insert default values
	team := "avengers"
	avengers := []karma.UserKarma{
		{User: "UCaptAmerica", Karma: 704},
		{User: "UHulk", Karma: 604},
		{User: "UThor", Karma: 1704},
		{User: "UDrStange", Karma: 505},
		{User: "UCaptMarvel", Karma: 9000},
		{User: "UIronman", Karma: 1400},
		{User: "USpiderman", Karma: 704},
	}
	for _, uk := range avengers {
		k, err := dao.UpdateKarma(team, uk.User, uk.Karma)
		assert.Nil(t, err)
		assert.Equal(t, uk.Karma, k)
	}

	// the names are hand sorted
	topNames := []string{"UCaptMarvel", "UThor", "UIronman", "USpiderman", "UCaptAmerica", "UHulk", "UDrStange"}

	// helper to convert users to names
	justNames := func(uk []karma.UserKarma) []string {
		names := make([]string, 0, len(uk))
		for _, u := range uk {
			names = append(names, u.User)
		}

		return names
	}

	topUsers, err := dao.Top(team, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(topUsers))
	assert.Equal(t, topNames[:1], justNames(topUsers))

	topUsers, err = dao.Top(team, 25)
	assert.Nil(t, err)
	assert.Equal(t, len(topNames), len(topUsers))
	assert.Equal(t, topNames, justNames(topUsers))

	topUsers, err = dao.Top("does-not-exist", 3)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(topUsers))
}
