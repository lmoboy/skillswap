package config

import (
	"os"
)

// SetupTestEnvironment configures the application for testing
func SetupTestEnvironment() {
	// Set test-specific environment variables
	os.Setenv("GO_ENV", "test")
	os.Setenv("SESSION_KEY", "test-session-key-for-testing-only")

	// Disable debug logging in tests
	os.Setenv("DEBUG", "false")
}
