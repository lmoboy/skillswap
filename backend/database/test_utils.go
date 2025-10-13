package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// TestDB is a test database connection
var TestDB *sql.DB

// SetupTestDB initializes a test database connection and creates test schema
func SetupTestDB() error {
	var err error
	TestDB, err = sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %v", err)
	}

	if err := TestDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping test database: %v", err)
	}

	// Create test tables
	if err := createTestTables(); err != nil {
		return fmt.Errorf("failed to create test tables: %v", err)
	}

	return nil
}

// createTestTables creates the necessary tables for testing
func createTestTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			profile_picture VARCHAR(500),
			aboutme TEXT,
			profession VARCHAR(255),
			location VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS skills (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) UNIQUE NOT NULL,
			description TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS user_skills (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			skill_id INT NOT NULL,
			verified TINYINT DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS courses (
			id INT AUTO_INCREMENT PRIMARY KEY,
			instructor_id INT NOT NULL,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10,2),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (instructor_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS course_assets (
			id INT AUTO_INCREMENT PRIMARY KEY,
			course_id INT NOT NULL,
			file_name VARCHAR(255) NOT NULL,
			file_path VARCHAR(500) NOT NULL,
			file_type VARCHAR(100) NOT NULL,
			file_size BIGINT NOT NULL,
			upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range tables {
		if _, err := TestDB.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
	}

	return nil
}

// ClearTestData removes all test data from the database
func ClearTestData() error {
	tables := []string{"user_skills", "course_assets", "courses", "users", "skills"}

	for _, table := range tables {
		if _, err := TestDB.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return fmt.Errorf("failed to clear table %s: %v", table, err)
		}
	}

	// Reset auto-increment counters
	for _, table := range tables {
		if _, err := TestDB.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", table)); err != nil {
			log.Printf("Warning: failed to reset auto-increment for %s: %v", table, err)
		}
	}

	return nil
}

// TeardownTestDB closes the test database connection
func TeardownTestDB() error {
	if TestDB != nil {
		return TestDB.Close()
	}
	return nil
}

// InsertTestUser creates a test user in the database
func InsertTestUser(username, email, password string) (int64, error) {
	result, err := TestDB.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES (?, ?, MD5(?))",
		username, email, password,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertTestSkill creates a test skill in the database
func InsertTestSkill(name, description string) (int64, error) {
	result, err := TestDB.Exec(
		"INSERT INTO skills (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
