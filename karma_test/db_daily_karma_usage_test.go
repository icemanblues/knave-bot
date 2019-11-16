package karma_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateKarmaDaily(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping db integration test")
	}

	db, dao, err := setupDB(testDB)
	assert.Nil(t, err)

	// confirm setup

	karmaRowCount := rowCountKarma(t, db)
	assert.Zero(t, karmaRowCount)

	dailyRowCount := rowCountDailyUsage(t, db)
	assert.Zero(t, dailyRowCount)

	date := time.Date(2019, time.November, 9, 0, 0, 0, 0, time.Local)

	// confirm update karma daily

	karmaSpiderman, err := dao.UpdateKarmaDaily("avengers", "ironman", "spiderman", 5, date)
	assert.Nil(t, err)
	assert.Equal(t, 5, karmaSpiderman)

	karmaRowCount = rowCountKarma(t, db)
	assert.Equal(t, 1, karmaRowCount)

	dailyRowCount = rowCountDailyUsage(t, db)
	assert.Equal(t, 1, dailyRowCount)

	usageIronman, err := dao.GetDaily("avengers", "ironman", date)
	assert.Nil(t, err)
	assert.Equal(t, 5, usageIronman)

	// confirm another update karma daily

	karmaSpiderman, err = dao.UpdateKarmaDaily("avengers", "ironman", "spiderman", -10, date)
	assert.Nil(t, err)
	assert.Equal(t, -5, karmaSpiderman)

	karmaRowCount = rowCountKarma(t, db)
	assert.Equal(t, 1, karmaRowCount)

	dailyRowCount = rowCountDailyUsage(t, db)
	assert.Equal(t, 1, dailyRowCount)

	usageIronman, err = dao.GetDaily("avengers", "ironman", date)
	assert.Nil(t, err)
	assert.Equal(t, 15, usageIronman)
}
