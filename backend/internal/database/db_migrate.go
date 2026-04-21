package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// splitSQL splits SQL content into individual statements
func splitSQL(content string) []string {
	var statements []string
	var current strings.Builder

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		if idx := strings.Index(line, "--"); idx != -1 {
			line = line[:idx]
		}

		current.WriteString(line)
		current.WriteString("\n")

		if strings.HasSuffix(strings.TrimSpace(line), ";") {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" && stmt != ";" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}

	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" && stmt != ";" {
			statements = append(statements, stmt)
		}
	}

	return statements
}

func translateMySQLToSQLite(sql string) string {
	// Remove MySQL engine/charset/collate
	reEngine := regexp.MustCompile(`(?i)ENGINE\s*=\s*\w+`)
	sql = reEngine.ReplaceAllString(sql, "")

	reCharset := regexp.MustCompile(`(?i)DEFAULT\s+CHARSET\s*=\s*\w+`)
	sql = reCharset.ReplaceAllString(sql, "")

	reCollate := regexp.MustCompile(`(?i)COLLATE\s*=\s*\w+`)
	sql = reCollate.ReplaceAllString(sql, "")

	// Fix DEFAULT("value") -> DEFAULT 'value'
	reDefaultParens := regexp.MustCompile(`(?i)DEFAULT\s*\(([^)]+)\)`)
	sql = reDefaultParens.ReplaceAllString(sql, "DEFAULT $1")

	// Translate BIGINT UNSIGNED NOT NULL AUTO_INCREMENT to INTEGER PRIMARY KEY
	rePK := regexp.MustCompile(`(?i)id\s+BIGINT\s+UNSIGNED\s+NOT\s+NULL\s+AUTO_INCREMENT`)
	sql = rePK.ReplaceAllString(sql, "id INTEGER PRIMARY KEY AUTOINCREMENT")

	// If it already has PRIMARY KEY (id), remove that since we handled it above
	rePKLine := regexp.MustCompile(`(?i),\s*PRIMARY\s+KEY\s*\(id\)`)
	sql = rePKLine.ReplaceAllString(sql, "")

	// Remove MySQL-specific key declarations
	reUniqueKey := regexp.MustCompile(`(?i)UNIQUE\s+KEY\s+\w+\s+\(`)
	sql = reUniqueKey.ReplaceAllString(sql, "UNIQUE (")

	// Remove stand-alone KEY lines and handle trailing commas
	reKey := regexp.MustCompile(`(?i),\s*KEY\s+\w+\s+\([^)]+\)`)
	sql = reKey.ReplaceAllString(sql, "")

	// Sometimes KEY is the last thing before ), so we need to handle that too
	reKeyLast := regexp.MustCompile(`(?i)KEY\s+\w+\s+\([^)]+\)\s*\)`)
	sql = reKeyLast.ReplaceAllString(sql, ")")
	// Remove "ON UPDATE CURRENT_TIMESTAMP"
	reOnUpdate := regexp.MustCompile(`(?i)ON\s+UPDATE\s+CURRENT_TIMESTAMP`)
	sql = reOnUpdate.ReplaceAllString(sql, "")

	// Replace ENUM with TEXT
	reEnum := regexp.MustCompile(`(?i)ENUM\s*\([^)]+\)`)
	sql = reEnum.ReplaceAllString(sql, "TEXT")

	// Replace INSERT IGNORE with INSERT OR IGNORE
	reInsertIgnore := regexp.MustCompile(`(?i)INSERT\s+IGNORE`)
	sql = reInsertIgnore.ReplaceAllString(sql, "INSERT OR IGNORE")

	// Replace SET FOREIGN_KEY_CHECKS
	reFKChecks := regexp.MustCompile(`(?i)SET\s+FOREIGN_KEY_CHECKS\s*=\s*[01];`)
	sql = reFKChecks.ReplaceAllString(sql, "")

	// Final cleanup: if we removed a key and left a trailing comma before a ), remove it
	reTrailingComma := regexp.MustCompile(`,\s*\)`)
	sql = reTrailingComma.ReplaceAllString(sql, ")")

	return sql
}
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			filename TEXT NOT NULL UNIQUE,
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
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	baseDir := filepath.Dir(execPath)

	migrationsDir := filepath.Join(baseDir, "migrations")
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		cwd, _ := os.Getwd()
		migrationsDir = filepath.Join(cwd, "internal", "database", "migrations")
		files, err = os.ReadDir(migrationsDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
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
		executed, err := isMigrationExecuted(db, filename)
		if err != nil {
			return err
		}
		if executed {
			continue
		}

		filePath := filepath.Join(migrationsDir, filename)
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		queries := splitSQL(string(content))
		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			translated := translateMySQLToSQLite(query)
			if translated == "" {
				continue
			}

			_, err = db.Exec(translated)
			if err != nil {
				// Log the error but keep going for some statements
				log.Printf("Warning: statement failed in migration %s: %v\nStatement: %s", filename, err, translated)
				// If it's a critical error (like a syntax error in a CREATE TABLE), we should still return it
				if strings.Contains(strings.ToLower(translated), "create table") {
					return err
				}
			}
		}

		if err := recordMigration(db, filename); err != nil {
			return err
		}
	}

	return nil
}
