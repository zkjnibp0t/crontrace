package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const schemaVersion = 1

const createTablesSQL = `
CREATE TABLE IF NOT EXISTS job_runs (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	command     TEXT    NOT NULL,
	started_at  DATETIME NOT NULL,
	finished_at DATETIME,
	duration_ms INTEGER,
	exit_code   INTEGER,
	stdout      TEXT,
	stderr      TEXT,
	created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_job_runs_command    ON job_runs(command);
CREATE INDEX IF NOT EXISTS idx_job_runs_started_at ON job_runs(started_at);
`

// Open opens (or creates) the SQLite database at the given path and
// ensures the schema is up to date.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(createTablesSQL)
	return err
}
