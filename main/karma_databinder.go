package main

import (
	"database/sql"
	"fmt"
)

type KarmaDB interface {
	GetKarma(team, user string) int
	UpdateKarma(team, user string, delta int) int
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
	return 0
}

func NewKdb(db *sql.DB) *LiteKarmaDB {
	return &LiteKarmaDB{db}
}
