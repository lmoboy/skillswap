package main

import (
	"os"
	"testing"

	"skillswap/backend/internal/config"
)

// TestMain sets up test environment and runs tests
func TestMain(m *testing.M) {
	// Set test environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_password")
	os.Setenv("DB_NAME", "skillswap_test")

	// Set up test configuration
	config.SetupTestEnvironment()

	// Run tests
	code := m.Run()

	// Cleanup
	os.Exit(code)
}
