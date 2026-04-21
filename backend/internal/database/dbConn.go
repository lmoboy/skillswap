package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// Init initializes the database connection and ensures the database and tables exist
func Init() {
	var err error

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		// Default to a local file if not set
		dbURL = "./skillswap.db"
	}

	// Ensure the directory for the database file exists
	dbDir := filepath.Dir(dbURL)
	if dbDir != "." {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}
	}

	// Connect to SQLite
	db, err = sql.Open("sqlite3", dbURL + "?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		if os.Getenv("GO_ENV") == "test" {
			panic(fmt.Sprintf("Failed to connect to SQLite: %v", err))
		}
		log.Fatal("Failed to connect to SQLite:", err)
	}

	// Ping to verify connection
	err = db.Ping()
	if err != nil {
		if os.Getenv("GO_ENV") == "test" {
			panic(fmt.Sprintf("Could not connect to SQLite: %v", err))
		}
		log.Fatal("Could not connect to SQLite:", err)
	}

	// Run migrations
	err = Migrate(db)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Verify critical tables exist
	if err := verifyTables(db); err != nil {
		log.Fatal("Database verification failed:", err)
	}
}

// verifyTables checks if critical tables exist
func verifyTables(db *sql.DB) error {
	criticalTables := []string{"users", "skills", "user_skills", "chats", "messages"}

	for _, table := range criticalTables {
		var exists bool
		query := fmt.Sprintf("SELECT COUNT(*) > 0 FROM sqlite_master WHERE type='table' AND name='%s'", table)
		err := db.QueryRow(query).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("critical table '%s' does not exist", table)
		}
	}

	return nil
}

func GetDatabase() (*sql.DB, error) {
	return db, nil
}

func Execute(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func Close() error {
	return db.Close()
}
