package karma

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/icemanblues/knave-bot/slack"
	log "github.com/sirupsen/logrus"
)

// UserKarma slack users and their karma totals
type UserKarma struct {
	User  string
	Karma int
}

// DAO Data Access Object for the Karma database
type DAO interface {
	GetKarma(team, user string) (int, error)
	UpdateKarma(team, user string, delta int) (int, error)
	DeleteKarma(team, user string) (int, error)
	Top(team string, n int) ([]UserKarma, error)
	Usage(slack.CommandData, slack.Response) error
	GetDaily(team, user string, date time.Time) (int, error)
	UpdateDaily(team, user string, date time.Time, karma int) (int, error)
}

// IsoDate converts a time object to 2006-01-02 format
func IsoDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// stringAttachment we use a *string so it can be nullable. This is to match the db column
func stringAttachment(att []slack.Attachments) *string {
	// it is null or empty, nothing to do here
	if att == nil || len(att) == 0 {
		return nil
	}

	// marshal it to json and save it
	j, err := json.Marshal(att)
	if err != nil {
		log.Warn("Unable to marshal attachments to json string. Inserting anyway without attachments", err)
		return nil
	}

	s := string(j)
	return &s
}

// SQLiteDAO a SQLite imlpementation of the Karma database
type SQLiteDAO struct {
	db *sql.DB
}

// GetKarma returns the karma value for the user in a given team
func (dao SQLiteDAO) GetKarma(team, user string) (int, error) {
	row := dao.db.QueryRow(`
		SELECT k.karma
		FROM   karma k
		WHERE  k.team = ?
		AND	   k.user = ?;
	`, team, user)

	var k int
	err := row.Scan(&k)
	if err != nil {
		log.Error("Unable to scan the row. It must be empty, query returned 0 rows.", team, user, err)
		return 0, nil
	}

	return k, nil
}

// UpdateKarma adds (or removes) karma from a user in a given team (workspace)
func (dao SQLiteDAO) UpdateKarma(workspace, user string, delta int) (int, error) {
	_, err := dao.db.Exec(`
		INSERT INTO karma
		(team, user, karma, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?)
		ON CONFLICT(team, user) DO UPDATE SET 
		karma = karma + excluded.karma,
		updated_at = excluded.updated_at;
	`, workspace, user, delta, time.Now(), time.Now())
	if err != nil {
		log.Error("Could not Insert or Update karma.", workspace, user, delta, err)
		return 0, err
	}

	return dao.GetKarma(workspace, user)
}

// DeleteKarma resets all karma for a given user in a given team to zer0
func (dao SQLiteDAO) DeleteKarma(team, user string) (int, error) {
	_, err := dao.db.Exec(`
		DELETE FROM karma
		WHERE  team = ?
		AND	   user = ?;
	`, team, user)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

// Usage tracks the usage of karma by pairing the request with the response
func (dao SQLiteDAO) Usage(data slack.CommandData, res slack.Response) error {
	s := stringAttachment(res.Attachments)

	_, err := dao.db.Exec(`
		INSERT INTO usage
		(command, text, enterprise, team, channel, user, created_at, response, response_type, attachments)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, data.Command, data.Text, data.EnterpriseID, data.TeamID, data.ChannelID, data.UserID, time.Now(), res.Text, res.ResponseType, s)

	return err
}

// Top returns the top n users (ordered by karma) from a given team
func (dao SQLiteDAO) Top(team string, n int) ([]UserKarma, error) {
	rows, err := dao.db.Query(`
		SELECT		k.user, 
					k.karma
		FROM  		karma k
		WHERE 		k.team = ?
		AND			k.karma >= 0
		ORDER BY	k.karma DESC, k.updated_at DESC
		LIMIT ?;
	`, team, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topUsers := make([]UserKarma, 0, n)
	for rows.Next() {
		var u UserKarma
		if err := rows.Scan(&u.User, &u.Karma); err != nil {
			log.Errorf("Unable to scan User Karma row for team %v %v", team, err)
		}
		topUsers = append(topUsers, u)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Errorf("Unable to scan User Karma from iterating rows for team %v", team)
		return nil, err
	}

	return topUsers, nil
}

// GetDaily return the amount of karma the team/user has gives/taken for a day
func (dao SQLiteDailyDAO) GetDaily(team, user string, date time.Time) (int, error) {
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
func (dao SQLiteDailyDAO) UpdateDaily(team, user string, date time.Time, karma int) (int, error) {
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


// NewDao factory method
func NewDao(db *sql.DB) SQLiteDAO {
	return SQLiteDAO{db}
}
