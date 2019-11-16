package karma_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func rowCountDailyUsage(t *testing.T, db *sql.DB) int {
	row := db.QueryRow("SELECT count(*) FROM daily_usage")
	var rowCount int
	err := row.Scan(&rowCount)
	assert.Nil(t, err)
	return rowCount
}

var date = time.Date(2019, time.November, 9, 0, 0, 0, 0, time.Local)

var noUsageDate = time.Date(2019, time.November, 22, 0, 0, 0, 0, time.Local)

func TestGetDailyLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	// confirm its empty
	rc := rowCountDailyUsage(t, db)
	assert.Zero(t, rc)

	usage, err := dao.GetDaily("yankees", "judge", date)
	assert.Nil(t, err)
	assert.Zero(t, usage)

	// insert one row
	kc, err := dao.UpdateDaily("yankees", "judge", date, 4)
	assert.Nil(t, err)
	assert.Equal(t, 4, kc)

	// confirm the row matches the dates
	rc = rowCountDailyUsage(t, db)
	assert.Equal(t, 1, rc)

	usage, err = dao.GetDaily("yankees", "judge", date)
	assert.Nil(t, err)
	assert.Equal(t, 4, usage)

	noUsage, err := dao.GetDaily("yankees", "judge", noUsageDate)
	assert.Nil(t, err)
	assert.Zero(t, noUsage)
}

func TestUpdateDailyLimit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	// confirm empty
	rc := rowCountDailyUsage(t, db)
	assert.Zero(t, rc)

	usage, err := dao.GetDaily("yankees", "judge", date)
	assert.Nil(t, err)
	assert.Zero(t, usage)

	// add and confirm the insert
	kc, err := dao.UpdateDaily("yankees", "judge", date, 4)
	assert.Nil(t, err)
	assert.Equal(t, 4, kc)

	usage, err = dao.GetDaily("yankees", "judge", date)
	assert.Nil(t, err)
	assert.Equal(t, 4, usage)

	rc = rowCountDailyUsage(t, db)
	assert.Equal(t, 1, rc)

	// add and confirm the update
	kc, err = dao.UpdateDaily("yankees", "judge", date, 9)
	assert.Nil(t, err)
	assert.Equal(t, 13, kc)

	usage, err = dao.GetDaily("yankees", "judge", date)
	assert.Nil(t, err)
	assert.Equal(t, 13, usage)

	// its an update, so the number of rows should remain at 1
	rc = rowCountDailyUsage(t, db)
	assert.Equal(t, 1, rc)
}
