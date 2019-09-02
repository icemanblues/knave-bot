package karma

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // sqlite db driver
)

// InitDB initializes the db and tables that we require
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := createDB(dataSourceName)
	if err != nil {
		return nil, err
	}

	err = schema(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDB(dataSourceName string) (*sql.DB, error) {
	// create the directories required for the database file
	dir := filepath.Dir(dataSourceName)
	err := os.MkdirAll(dir, os.ModeDir|0755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func schema(db *sql.DB) error {
	// karma table
	if err := schemaKarma(db); err != nil {
		return err
	}

	// usage table
	if err := schemaUsage(db); err != nil {
		return err
	}

	return nil
}

func schemaKarma(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS karma (
		id			INTEGER PRIMARY KEY, 
		team		TEXT, 
		user		TEXT,
		karma 		INTEGER,
		created_at	TEXT,
		updated_at	TEXT
	);
	CREATE UNIQUE INDEX IF NOT EXISTS idx_karma_team_user ON karma (team, user);
	`)

	return err
}

func schemaUsage(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS usage (
		command			TEXT,
		text			TEXT,
		enterprise		TEXT,
		team			TEXT,
		channel			TEXT,
		user			TEXT,
		created_at		TEXT,
		response		TEXT,
		response_type	TEXT,
		attachments 	TEXT
	);
	`)

	return err
}
