package karma

import (
	"database/sql"
	"time"
)

// DAO Data Access Object for the Karma database
type DAO interface {
	GetKarma(team, user string) (int, error)
	UpdateKarma(team, user string, delta int) (int, error)
	DeleteKarma(team, user string) (int, error)
}

// SQLiteDAO a SQLite imlpementation of the Karma database
type SQLiteDAO struct {
	db *sql.DB
}

// GetKarma returns the karma value for the user in a given team
func (kdb *SQLiteDAO) GetKarma(team, user string) (int, error) {
	row := kdb.db.QueryRow(`
		SELECT k.karma
		FROM   karma k
		WHERE  k.team = ?
		AND	   k.user = ?;
	`, team, user)

	var k int
	err := row.Scan(&k)
	if err != nil {
		return 0, nil
	}

	return k, nil
}

// UpdateKarma adds (or removes) karma from a user in a given team (workspace)
func (kdb *SQLiteDAO) UpdateKarma(workspace, user string, delta int) (int, error) {
	_, err := kdb.db.Exec(`
		INSERT INTO karma
		(team, user, karma, created_at, updated_at)
		VALUES
		(?, ?, ?, ?, ?)
		ON CONFLICT(team, user) DO UPDATE SET 
		karma = karma + excluded.karma,
		updated_at = excluded.updated_at;
	`, workspace, user, delta, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	return kdb.GetKarma(workspace, user)
}

// DeleteKarma resets all karma for a given user in a given team to zer0
func (kdb *SQLiteDAO) DeleteKarma(team, user string) (int, error) {
	_, err := kdb.db.Exec(`
		DELETE FROM karma
		WHERE  team = ?
		AND	   user = ?;
	`, team, user)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

// NewKdb factory method
func NewKdb(db *sql.DB) *SQLiteDAO {
	return &SQLiteDAO{db}
}
