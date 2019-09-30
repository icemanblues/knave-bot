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
	Usage(*slack.CommandData, *slack.Response) error
	Top(team string, n int) ([]UserKarma, error)
}

// SQLiteDAO a SQLite imlpementation of the Karma database
type SQLiteDAO struct {
	db *sql.DB
}

// GetKarma returns the karma value for the user in a given team
func (dao *SQLiteDAO) GetKarma(team, user string) (int, error) {
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
func (dao *SQLiteDAO) UpdateKarma(workspace, user string, delta int) (int, error) {
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
func (dao *SQLiteDAO) DeleteKarma(team, user string) (int, error) {
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

// Usage tracks the usage of karma by pairing the request with the response
func (dao *SQLiteDAO) Usage(data *slack.CommandData, res *slack.Response) error {
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
func (dao *SQLiteDAO) Top(team string, n int) ([]UserKarma, error) {
	rows, err := dao.db.Query(`
		SELECT		k.user, 
					k.karma
		FROM  		karma k
		WHERE 		k.team = ?
		AND			k.karma > 0
		ORDER BY	k.karma DESC
		LIMIT ?;
	`, team, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	topUser := make([]UserKarma, 0, n)
	for rows.Next() {
		var u UserKarma
		if err := rows.Scan(&u.User, &u.Karma); err != nil {
			log.Errorf("Unable to scan User Karma row for team %v", team)
		}
		topUser = append(topUser, u)
	}
	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		log.Errorf("Unable to scan User Karma from iterating rows for team %v", team)
		return nil, err
	}

	return topUser, nil
}

// NewDao factory method
func NewDao(db *sql.DB) *SQLiteDAO {
	return &SQLiteDAO{db}
}
