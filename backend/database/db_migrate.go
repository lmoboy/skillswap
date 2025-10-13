package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"skillswap/backend/utils"
	"sort"
	"strings"
)

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			filename VARCHAR(255) NOT NULL UNIQUE,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func recordMigration(db *sql.DB, filename string) error {
	_, err := db.Exec("INSERT INTO migrations (filename) VALUES (?)", filename)
	return err
}

func Migrate(db *sql.DB) error {
	// Create migrations tracking table
	if err := createMigrationsTable(db); err != nil {
		utils.DebugPrint("failed to create migrations table: %v", err)
		return err
	}

	// Get the executable's directory or current working directory
	execPath, err := os.Executable()
	if err != nil {
		utils.DebugPrint("failed to get executable path: %v", err)
		return err
	}
	baseDir := filepath.Dir(execPath)

	// Try to find migrations directory
	migrationsDir := filepath.Join(baseDir, "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		// Fallback to checking current directory
		cwd, _ := os.Getwd()
		migrationsDir = filepath.Join(cwd, "database", "migrations")
		files, err = os.ReadDir(migrationsDir)
		if err != nil {
			if os.IsNotExist(err) {
				utils.DebugPrint("Migrations directory not found at:", migrationsDir)
				return nil
			}
			utils.DebugPrint("failed to read migrations directory: %v", err)
			return err
		}
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		// Check if migration already executed

		filePath := filepath.Join(migrationsDir, filename)
		utils.DebugPrint("Executing migration:", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			utils.DebugPrint("failed to read migration file %s: %v", filename, err)
			return err
		}

		queries := strings.Split(string(content), ";")
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" || strings.HasPrefix(query, "--") {
				continue
			}
			_, err = db.Exec(query)
			if err != nil {
				utils.DebugPrint("failed to execute migration %s: %v", filename, err)
				return err
			}
		}

		utils.DebugPrint("Successfully executed migration:", filename)
	}

	return nil
}
