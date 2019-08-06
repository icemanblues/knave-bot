package main

import (
	"database/sql"
	"fmt"
	"time"
)

type KarmaDB interface {
	GetKarma(team, user string) int
	UpdateKarma(team, user string, delta int) int
	DeleteKarma(team, user string) int
}

type LiteKarmaDB struct {
	db *sql.DB
}

func (kdb *LiteKarmaDB) GetKarma(team, user string) int {
	row := kdb.db.QueryRow(`
		SELECT k.karma
		FROM   karma k
		WHERE  k.team = ?
		AND	   k.user = ?
	`, team, user)

	var k int
	err := row.Scan(&k)
	if err != nil {
		fmt.Printf("WARN: Didn't find the user %v in team %v. Defaulting to %v error: %v\n", user, team, k, err)
	}

	return k
}

func (kdb *LiteKarmaDB) UpdateKarma(workspace, user string, delta int) int {
	// TODO: UpSert the current value plus delta. If no current value, assume 0
	// currently we do a SELECT to determine if we need to do an insert or update
	// in the future, we will combine this into a proper UpSert
	// and use the returning syntax to make it atomic
	k := kdb.GetKarma(workspace, user)
	if k == 0 {
		// insert
		_, err := kdb.db.Exec(`
			INSERT INTO karma
			(team, user, karma, created_at, updated_at)
			VALUES
			(?, ?, ?, ?, ?)
		`, workspace, user, delta, time.Now(), time.Now())
		if err != nil {
			fmt.Printf("ERROR: Unable to insert a new value to the karma %v %v. err %v\n", workspace, user, err)
			return 0
		}
		return delta
	}

	// update
	_, err := kdb.db.Exec(`
		UPDATE karma
		SET	karma = karma + ?,
			updated_at = ?
		WHERE  team = ?
		AND	   user = ? 
	`, delta, time.Now(), workspace, user)
	if err != nil {
		fmt.Printf("ERROR: Unable to update karma. delta %v workspace %v user %v. err %v\n", delta, workspace, user, err)
		return k
	}
	return kdb.GetKarma(workspace, user)
}

func (kdb *LiteKarmaDB) DeleteKarma(team, user string) int {
	_, err := kdb.db.Exec(`
		DELETE FROM karma
		WHERE  team = ?
		AND	   user = ? 
	`, team, user)
	if err != nil {
		fmt.Printf("ERROR: Unable to delete karma for team %v user %v. err %v\n", team, user, err)
	}

	return 0
}

func NewKdb(db *sql.DB) *LiteKarmaDB {
	return &LiteKarmaDB{db}
}
