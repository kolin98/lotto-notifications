package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func GetDB() (*sqlx.DB, error) {
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return db, nil
}

// Initialize sets up the database connection and creates the database if it doesn't exist
func Initialize(dbPath string) error {
	// Ensure data directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database connection using sqlx
	dbConnection, err := sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	dbConnection.SetMaxOpenConns(1) // SQLite only supports one writer at a time
	dbConnection.SetMaxIdleConns(1)

	// Configure SQLite settings
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = normal",
		"PRAGMA journal_size_limit = 6144000",
	}

	for _, pragma := range pragmas {
		if _, err := dbConnection.Exec(pragma); err != nil {
			dbConnection.Close()
			return fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	db = dbConnection

	return nil
}

// Close closes the database connection
func Close() error {
	if db != nil {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
	}
	return nil
}
