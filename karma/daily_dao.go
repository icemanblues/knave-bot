package karma

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

// DailyDao provides the DAO interface for checking and UPdating daily limits
type DailyDao interface {
	GetDaily(team, user string, date time.Time) (int, error)
	UpdateDaily(team, user string, date time.Time, karma int) (int, error)
}

// SQLiteDailyDAO .
type SQLiteDailyDAO struct {
	db *sql.DB
}

const layout string = "2006-01-02"

// IsoDate converts a time object to 2006-01-02 format
func IsoDate(t time.Time) string {
	return t.Format(layout)
}

// GetDaily return the amount of karma the team/user has gives/taken for a day
func (dao *SQLiteDailyDAO) GetDaily(team, user string, date time.Time) (int, error) {
	row := dao.db.QueryRow(`
		SELECT du.usage
		FROM   daily_usage du
		WHERE  du.team = ?
		AND	   du.user = ?
		AND	   du.daily = ?;
	`, team, user, IsoDate(date))

	var u int
	err := row.Scan(&u)
	if err != nil {
		log.Error("Unable to scan the row. It must be empty, query returned 0 rows.", team, user, err)
		return 0, nil
	}

	return u, nil
}

// UpdateDaily adds karma to team/user's daily usage count
func (dao *SQLiteDailyDAO) UpdateDaily(team, user string, date time.Time, karma int) (int, error) {
	_, err := dao.db.Exec(`
		INSERT INTO daily_usage
		(team, user, daily, usage, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?, ?)
		ON CONFLICT(team, user, daily) DO UPDATE SET 
		usage = usage + excluded.usage,
		updated_at = excluded.updated_at;
	`, team, user, IsoDate(date), karma, time.Now(), time.Now())
	if err != nil {
		log.Error("Could not Insert or Update karma.", team, user, date, karma, err)
		return 0, err
	}

	return dao.GetDaily(team, user, date)
}

// NewDailyDao factory method
func NewDailyDao(db *sql.DB) *SQLiteDailyDAO {
	return &SQLiteDailyDAO{db}
}
