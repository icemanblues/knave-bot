package karma

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // sqlite db driver
)

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
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS karma (
			id INTEGER PRIMARY KEY, 
			team TEXT, 
			user TEXT,
			karma INTEGER,
			created_at TEXT,
			updated_at TEXT
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_karma_team_user ON karma (team, user);
	`)
	if err != nil {
		return err
	}

	return nil
}

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
