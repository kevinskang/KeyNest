package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Open opens the SQLite database at path, applies required PRAGMAs, and runs migrations.
func Open(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	// SQLite performs best with a single connection; also ensures PRAGMAs apply globally.
	db.SetMaxOpenConns(1)

	if _, err = db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		db.Close()
		return nil, fmt.Errorf("set foreign_keys pragma: %w", err)
	}
	if _, err = db.Exec(`PRAGMA journal_mode = WAL`); err != nil {
		db.Close()
		return nil, fmt.Errorf("set journal_mode pragma: %w", err)
	}

	if err = migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS schema_version (
			version    INTEGER PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id            INTEGER  PRIMARY KEY AUTOINCREMENT,
			email         TEXT     NOT NULL UNIQUE,
			password_hash TEXT     NOT NULL,
			enc_salt      TEXT     NOT NULL,
			created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		`CREATE TABLE IF NOT EXISTS api_keys (
			id              INTEGER  PRIMARY KEY AUTOINCREMENT,
			user_id         INTEGER  NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			key_name        TEXT     NOT NULL,
			key_value       TEXT     NOT NULL,
			url             TEXT,
			expiry_date     DATE,
			registered_date DATE,
			memo            TEXT,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_user_id    ON api_keys(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_key_name   ON api_keys(user_id, key_name)`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_expiry     ON api_keys(user_id, expiry_date)`,
		`CREATE INDEX IF NOT EXISTS idx_api_keys_created_at ON api_keys(user_id, created_at)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("migration statement failed: %w", err)
		}
	}
	return nil
}
