package warehouse

import (
	"os"
	"testing"
)

// TestMain sets up any global test configuration
func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()
	os.Exit(code)
}
