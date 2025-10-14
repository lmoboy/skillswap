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

func isMigrationExecuted(db *sql.DB, filename string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = ?", filename).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
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

	// Determine migration directory
	execPath, err := os.Executable()
	if err != nil {
		utils.DebugPrint("failed to get executable path: %v", err)
		return err
	}
	baseDir := filepath.Dir(execPath)

	migrationsDir := filepath.Join(baseDir, "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		// Fallback: check working directory
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

	// Sort and process SQL files
	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, filename := range sqlFiles {
		// Check if migration has already been applied
		executed, err := isMigrationExecuted(db, filename)
		if err != nil {
			utils.DebugPrint("failed to check migration %s: %v", filename, err)
			return err
		}
		if executed {
			utils.DebugPrint("Skipping already executed migration:", filename)
			continue
		}

		filePath := filepath.Join(migrationsDir, filename)
		utils.DebugPrint("Executing migration:", filePath)
		content, err := os.ReadFile(filePath)
		if err != nil {
			utils.DebugPrint("failed to read migration file %s: %v", filename, err)
			return err
		}

		// Split by semicolon and execute each query
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

		// Record successful migration
		if err := recordMigration(db, filename); err != nil {
			utils.DebugPrint("failed to record migration %s: %v", filename, err)
			return err
		}

		utils.DebugPrint("Successfully executed migration:", filename)
	}

	return nil
}
