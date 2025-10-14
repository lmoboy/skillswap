package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"skillswap/backend/utils"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Init initializes the database connection and ensures the database and tables exist
func Init() {
	var err error
	
	// First, connect to MySQL without specifying a database to check if skillswap DB exists
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Parse the DB_URL to extract the database name
	dbName := extractDatabaseName(dbURL)
	if dbName == "" {
		log.Fatal("Could not extract database name from DB_URL")
	}

	utils.DebugPrint("Connecting to MySQL server...")
	
	// Connect to MySQL server (without database)
	serverURL := strings.Replace(dbURL, "/"+dbName, "/", 1)
	serverDB, err := sql.Open("mysql", serverURL)
	if err != nil {
		log.Fatal("Failed to connect to MySQL server:", err)
	}

	// Ping to verify connection
	for i := 0; i < 10; i++ {
		err = serverDB.Ping()
		if err == nil {
			break
		}
		utils.DebugPrint(fmt.Sprintf("Waiting for MySQL server (attempt %d/10)...", i+1))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Could not connect to MySQL server after 10 attempts:", err)
	}

	utils.DebugPrint("Connected to MySQL server")

	// Create database if it doesn't exist
	_, err = serverDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}
	utils.DebugPrint(fmt.Sprintf("Database '%s' ready", dbName))

	serverDB.Close()

	// Now connect to the specific database
	db, err = sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	
	utils.DebugPrint("Database connected successfully")

	// Check if migrations table exists
	var tableExists bool
	err = db.QueryRow("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = ? AND table_name = 'migrations'", dbName).Scan(&tableExists)
	if err != nil {
		utils.DebugPrint("Error checking for migrations table:", err.Error())
		tableExists = false
	}

	// If migrations table doesn't exist, run migrations
	if !tableExists {
		utils.DebugPrint("Migrations table not found. Running migrations...")
		err = Migrate(db)
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
		utils.DebugPrint("Migrations completed successfully")
	} else {
		utils.DebugPrint("Migrations table found. Checking for pending migrations...")
		// Check if users table exists as a basic sanity check
		var usersExists bool
		err = db.QueryRow("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = ? AND table_name = 'users'", dbName).Scan(&usersExists)
		if err == nil && !usersExists {
			utils.DebugPrint("Users table not found. Running migrations...")
			err = Migrate(db)
			if err != nil {
				log.Fatal("Migration failed:", err)
			}
		} else {
			// Run migrations to check for any new ones
			err = Migrate(db)
			if err != nil {
				utils.DebugPrint("Warning: Migration check failed:", err.Error())
			} else {
				utils.DebugPrint("All migrations up to date")
			}
		}
	}

	// Verify critical tables exist
	if err := verifyTables(db, dbName); err != nil {
		log.Fatal("Database verification failed:", err)
	}

	utils.DebugPrint("Database initialization complete")
}

// extractDatabaseName extracts the database name from the connection string
func extractDatabaseName(dbURL string) string {
	// Format: user:pass@tcp(host:port)/database?params
	parts := strings.Split(dbURL, "/")
	if len(parts) < 2 {
		return ""
	}
	dbPart := parts[len(parts)-1]
	// Remove query parameters if present
	if idx := strings.Index(dbPart, "?"); idx != -1 {
		dbPart = dbPart[:idx]
	}
	return dbPart
}

// verifyTables checks if critical tables exist
func verifyTables(db *sql.DB, dbName string) error {
	criticalTables := []string{"users", "skills", "user_skills", "chats", "messages"}
	
	for _, table := range criticalTables {
		var exists bool
		err := db.QueryRow("SELECT COUNT(*) > 0 FROM information_schema.tables WHERE table_schema = ? AND table_name = ?", dbName, table).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}
		if !exists {
			return fmt.Errorf("critical table '%s' does not exist", table)
		}
	}
	
	utils.DebugPrint("All critical tables verified")
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

func Debug(query string, args ...interface{}) {
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logQueryError(err error) {
	if err != nil {
		log.Println("Query error:", err)
	}
}

func logQueryResult(result sql.Result, err error) {
	if err != nil {
		log.Println("Query result:", result, err)
	}
}

func logQueryRow(row *sql.Row) {
	if row == nil {
		log.Println("Query row: nil")
	} else {
		log.Println("Query row:", row)
	}
}
