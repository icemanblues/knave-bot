package karma

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/icemanblues/knave-bot/slack"
	log "github.com/sirupsen/logrus"
)

// DAO Data Access Object for the Karma database
type DAO interface {
	GetKarma(team, user string) (int, error)
	UpdateKarma(team, user string, delta int) (int, error)
	DeleteKarma(team, user string) (int, error)
	Usage(*slack.CommandData, *slack.Response) error
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

// Usage tracks the usage of karma by pairing the request with the response
func (dao *SQLiteDAO) Usage(data *slack.CommandData, res *slack.Response) error {
	var att *string
	j, err := json.Marshal(res.Attachments)
	if err != nil {
		log.Warn("Unable to json marshal attachments to string. Inserting anyway")
		*att = ("Error " + err.Error())
	} else {
		*att = string(j)
	}

	_, err = dao.db.Exec(`
		INSERT INTO usage
		(command, text, enterprise, team, channel, user, created_at, response, response_type, attachements)
		VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, data.Command, data.Text, data.EnterpriseID, data.TeamID, data.ChannelID, data.UserID, time.Now(), res.Text, res.ResponseType, att)

	return err
}

// NewDao factory method
func NewDao(db *sql.DB) *SQLiteDAO {
	return &SQLiteDAO{db}
}
