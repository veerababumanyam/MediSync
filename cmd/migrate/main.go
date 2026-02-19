// Package main provides a CLI tool for running PostgreSQL database migrations.
// This tool supports:
// - Running all pending migrations (up)
// - Rolling back migrations (down)
// - Migrating to a specific version
// - Showing current migration status
// - Creating new migration files
//
// Usage:
//
//	go run cmd/migrate/main.go [command] [options]
//
// Commands:
//
//	up        Run all pending migrations
//	down      Roll back the last migration
//	down-all  Roll back all migrations
//	version   Show current migration version
//	force     Force set the migration version (use with caution)
//	create    Create new migration files
//	status    Show migration status
//
// Environment Variables:
//
//	DATABASE_URL  PostgreSQL connection string (required)
//	              Format: postgres://user:password@host:port/dbname?sslmode=disable
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	// Default migrations directory relative to project root
	defaultMigrationsPath = "migrations"

	// Exit codes
	exitOK          = 0
	exitError       = 1
	exitUsageError  = 2
	exitNoChange    = 3 // No migrations to run
	exitDirty       = 4 // Database is in dirty state
)

// Config holds the migration tool configuration
type Config struct {
	DatabaseURL    string
	MigrationsPath string
	Verbose        bool
}

func main() {
	// Initialize structured logger
	logLevel := slog.LevelInfo
	if os.Getenv("LOG_LEVEL") == "debug" {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Parse command line arguments
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(exitUsageError)
	}

	// Setup flags for all commands
	flagSet := flag.NewFlagSet("migrate", flag.ExitOnError)
	databaseURL := flagSet.String("database-url", "", "PostgreSQL connection string (overrides DATABASE_URL env)")
	migrationsPath := flagSet.String("path", defaultMigrationsPath, "Path to migrations directory")
	verbose := flagSet.Bool("verbose", false, "Enable verbose output")

	// Find the command (first non-flag argument)
	command := ""
	commandIdx := 1

	// Look for the command among arguments
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if !strings.HasPrefix(arg, "-") {
			command = arg
			commandIdx = i
			break
		}
		// Skip the value of flags that take arguments
		if arg == "--database-url" || arg == "--path" || arg == "-database-url" || arg == "-path" {
			i++ // Skip the next argument (flag value)
		}
	}

	if command == "" {
		printUsage()
		os.Exit(exitUsageError)
	}

	// Build args slice: everything except the command
	var argsToparse []string
	for i := 1; i < len(os.Args); i++ {
		if i == commandIdx {
			continue
		}
		argsToparse = append(argsToparse, os.Args[i])
	}

	// Parse flags
	if err := flagSet.Parse(argsToparse); err != nil {
		slog.Error("Failed to parse flags", "error", err)
		os.Exit(exitUsageError)
	}

	// Build configuration
	config := Config{
		DatabaseURL:    *databaseURL,
		MigrationsPath: *migrationsPath,
		Verbose:        *verbose,
	}

	// Get database URL from environment if not provided via flag
	if config.DatabaseURL == "" {
		config.DatabaseURL = os.Getenv("DATABASE_URL")
	}

	// Handle commands that don't require database connection
	switch command {
	case "help", "-h", "--help":
		printUsage()
		os.Exit(exitOK)
	case "create":
		handleCreate(flagSet.Args(), config)
		return
	}

	// Validate database URL for commands that need it
	if config.DatabaseURL == "" {
		slog.Error("DATABASE_URL is required",
			"hint", "Set DATABASE_URL environment variable or use --database-url flag")
		os.Exit(exitUsageError)
	}

	// Convert relative path to absolute
	absPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		slog.Error("Failed to resolve migrations path", "path", config.MigrationsPath, "error", err)
		os.Exit(exitError)
	}
	config.MigrationsPath = absPath

	// Verify migrations directory exists
	if _, err := os.Stat(config.MigrationsPath); os.IsNotExist(err) {
		slog.Error("Migrations directory not found", "path", config.MigrationsPath)
		os.Exit(exitError)
	}

	// Create migrate instance
	sourceURL := fmt.Sprintf("file://%s", config.MigrationsPath)
	m, err := migrate.New(sourceURL, config.DatabaseURL)
	if err != nil {
		slog.Error("Failed to create migrate instance",
			"error", err,
			"source", sourceURL,
		)
		os.Exit(exitError)
	}
	defer m.Close()

	// Enable verbose logging if requested
	if config.Verbose {
		m.Log = &migrateLogger{logger: logger}
	}

	// Execute command
	switch command {
	case "up":
		handleUp(m, flagSet.Args())
	case "down":
		handleDown(m, flagSet.Args())
	case "down-all":
		handleDownAll(m)
	case "version":
		handleVersion(m)
	case "force":
		handleForce(m, flagSet.Args())
	case "status":
		handleStatus(m, config)
	case "drop":
		handleDrop(m)
	default:
		slog.Error("Unknown command", "command", command)
		printUsage()
		os.Exit(exitUsageError)
	}
}

// handleUp runs all pending migrations or migrates to a specific version
func handleUp(m *migrate.Migrate, args []string) {
	var err error

	if len(args) > 0 {
		// Migrate to specific version
		version, parseErr := strconv.ParseUint(args[0], 10, 64)
		if parseErr != nil {
			slog.Error("Invalid version number", "version", args[0], "error", parseErr)
			os.Exit(exitUsageError)
		}
		slog.Info("Migrating to version", "target_version", version)
		err = m.Migrate(uint(version))
	} else {
		// Run all pending migrations
		slog.Info("Running all pending migrations")
		err = m.Up()
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No migrations to run - database is up to date")
			os.Exit(exitNoChange)
		}
		slog.Error("Migration failed", "error", err)
		os.Exit(exitError)
	}

	version, dirty, _ := m.Version()
	slog.Info("Migration completed successfully",
		"current_version", version,
		"dirty", dirty,
	)
	os.Exit(exitOK)
}

// handleDown rolls back the last migration or a specific number of migrations
func handleDown(m *migrate.Migrate, args []string) {
	steps := 1
	if len(args) > 0 {
		var err error
		steps, err = strconv.Atoi(args[0])
		if err != nil || steps < 1 {
			slog.Error("Invalid step count", "steps", args[0])
			os.Exit(exitUsageError)
		}
	}

	slog.Info("Rolling back migrations", "steps", steps)
	err := m.Steps(-steps)

	if err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No migrations to roll back")
			os.Exit(exitNoChange)
		}
		slog.Error("Rollback failed", "error", err)
		os.Exit(exitError)
	}

	version, dirty, _ := m.Version()
	slog.Info("Rollback completed successfully",
		"current_version", version,
		"dirty", dirty,
	)
	os.Exit(exitOK)
}

// handleDownAll rolls back all migrations
func handleDownAll(m *migrate.Migrate) {
	slog.Warn("Rolling back ALL migrations - this will destroy all data!")

	err := m.Down()
	if err != nil {
		if err == migrate.ErrNoChange {
			slog.Info("No migrations to roll back - database is clean")
			os.Exit(exitNoChange)
		}
		slog.Error("Rollback failed", "error", err)
		os.Exit(exitError)
	}

	slog.Info("All migrations rolled back successfully")
	os.Exit(exitOK)
}

// handleVersion displays the current migration version
func handleVersion(m *migrate.Migrate) {
	version, dirty, err := m.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			slog.Info("No migrations have been applied",
				"version", "none",
				"dirty", false,
			)
			fmt.Println("Version: none (no migrations applied)")
			os.Exit(exitOK)
		}
		slog.Error("Failed to get version", "error", err)
		os.Exit(exitError)
	}

	slog.Info("Current migration version",
		"version", version,
		"dirty", dirty,
	)
	fmt.Printf("Version: %d\n", version)
	if dirty {
		fmt.Println("Status: DIRTY (migration was interrupted)")
		slog.Warn("Database is in dirty state - run 'force' command to fix")
		os.Exit(exitDirty)
	}
	fmt.Println("Status: clean")
	os.Exit(exitOK)
}

// handleForce sets the migration version without running migrations
func handleForce(m *migrate.Migrate, args []string) {
	if len(args) < 1 {
		slog.Error("Version number required for force command")
		fmt.Println("Usage: migrate force <version>")
		fmt.Println("  Use -1 to mark database as having no migrations")
		os.Exit(exitUsageError)
	}

	version, err := strconv.Atoi(args[0])
	if err != nil {
		slog.Error("Invalid version number", "version", args[0], "error", err)
		os.Exit(exitUsageError)
	}

	slog.Warn("Forcing migration version",
		"version", version,
		"warning", "This does not run any migrations, only sets the version marker",
	)

	if err := m.Force(version); err != nil {
		slog.Error("Failed to force version", "error", err)
		os.Exit(exitError)
	}

	slog.Info("Version forced successfully", "version", version)
	os.Exit(exitOK)
}

// handleStatus shows detailed migration status
func handleStatus(m *migrate.Migrate, config Config) {
	version, dirty, err := m.Version()

	fmt.Println("=== MediSync Migration Status ===")
	fmt.Printf("Migrations Path: %s\n", config.MigrationsPath)
	fmt.Println()

	if err != nil {
		if err == migrate.ErrNilVersion {
			fmt.Println("Database Version: none (no migrations applied)")
		} else {
			slog.Error("Failed to get version", "error", err)
			os.Exit(exitError)
		}
	} else {
		fmt.Printf("Database Version: %d\n", version)
		if dirty {
			fmt.Println("Status: DIRTY ⚠️  (migration was interrupted)")
			fmt.Println("  Run 'migrate force <version>' to fix the dirty state")
		} else {
			fmt.Println("Status: Clean ✓")
		}
	}

	// List available migration files
	fmt.Println()
	fmt.Println("Available Migrations:")

	files, err := os.ReadDir(config.MigrationsPath)
	if err != nil {
		slog.Error("Failed to read migrations directory", "error", err)
		os.Exit(exitError)
	}

	migrations := make(map[uint64]struct {
		up   bool
		down bool
		name string
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		// Parse migration filename: {version}_{name}.{up|down}.sql
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			continue
		}

		versionStr := parts[0]
		v, err := strconv.ParseUint(versionStr, 10, 64)
		if err != nil {
			continue
		}

		entry := migrations[v]

		// Extract migration name (everything between version and direction)
		nameEnd := strings.LastIndex(name, ".up.sql")
		if nameEnd == -1 {
			nameEnd = strings.LastIndex(name, ".down.sql")
		}
		if nameEnd > 0 {
			entry.name = name[len(versionStr)+1 : nameEnd]
		}

		if strings.HasSuffix(name, ".up.sql") {
			entry.up = true
		} else if strings.HasSuffix(name, ".down.sql") {
			entry.down = true
		}

		migrations[v] = entry
	}

	// Print migrations in order
	if len(migrations) == 0 {
		fmt.Println("  No migration files found")
	} else {
		for v := uint64(1); v <= uint64(len(migrations)+10); v++ {
			entry, ok := migrations[v]
			if !ok {
				continue
			}

			status := "pending"
			if version > 0 && v <= uint64(version) {
				status = "applied"
			}

			upIcon := "✗"
			if entry.up {
				upIcon = "✓"
			}
			downIcon := "✗"
			if entry.down {
				downIcon = "✓"
			}

			fmt.Printf("  %03d: %s [up:%s down:%s] - %s\n", v, entry.name, upIcon, downIcon, status)
		}
	}

	os.Exit(exitOK)
}

// handleCreate creates new migration files
func handleCreate(args []string, config Config) {
	if len(args) < 1 {
		slog.Error("Migration name required")
		fmt.Println("Usage: migrate create <name>")
		fmt.Println("Example: migrate create add_users_table")
		os.Exit(exitUsageError)
	}

	name := args[0]
	// Sanitize name - replace spaces with underscores, remove special chars
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	// Find next version number
	absPath, err := filepath.Abs(config.MigrationsPath)
	if err != nil {
		slog.Error("Failed to resolve migrations path", "error", err)
		os.Exit(exitError)
	}

	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll(absPath, 0755); err != nil {
		slog.Error("Failed to create migrations directory", "error", err)
		os.Exit(exitError)
	}

	files, err := os.ReadDir(absPath)
	if err != nil {
		slog.Error("Failed to read migrations directory", "error", err)
		os.Exit(exitError)
	}

	maxVersion := uint64(0)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fname := file.Name()
		parts := strings.Split(fname, "_")
		if len(parts) >= 1 {
			if v, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
				if v > maxVersion {
					maxVersion = v
				}
			}
		}
	}

	nextVersion := maxVersion + 1
	timestamp := time.Now().Format("20060102150405")

	// Create migration files with timestamp for uniqueness
	upFile := filepath.Join(absPath, fmt.Sprintf("%03d_%s.up.sql", nextVersion, name))
	downFile := filepath.Join(absPath, fmt.Sprintf("%03d_%s.down.sql", nextVersion, name))

	upContent := fmt.Sprintf(`-- Migration: %s
-- Created: %s
-- Version: %d

-- Write your UP migration SQL here

`, name, timestamp, nextVersion)

	downContent := fmt.Sprintf(`-- Migration: %s (rollback)
-- Created: %s
-- Version: %d

-- Write your DOWN migration SQL here
-- This should undo the changes made in the UP migration

`, name, timestamp, nextVersion)

	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		slog.Error("Failed to create up migration", "file", upFile, "error", err)
		os.Exit(exitError)
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		slog.Error("Failed to create down migration", "file", downFile, "error", err)
		os.Exit(exitError)
	}

	slog.Info("Created migration files",
		"version", nextVersion,
		"name", name,
		"up_file", upFile,
		"down_file", downFile,
	)

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  Up:   %s\n", upFile)
	fmt.Printf("  Down: %s\n", downFile)

	os.Exit(exitOK)
}

// handleDrop drops everything in the database (dangerous!)
func handleDrop(m *migrate.Migrate) {
	slog.Warn("Dropping all database objects - THIS IS DESTRUCTIVE!")
	fmt.Println("WARNING: This will drop all tables and data!")
	fmt.Println("Press Ctrl+C to cancel, or wait 5 seconds to continue...")

	time.Sleep(5 * time.Second)

	if err := m.Drop(); err != nil {
		slog.Error("Failed to drop database", "error", err)
		os.Exit(exitError)
	}

	slog.Info("Database dropped successfully")
	os.Exit(exitOK)
}

// migrateLogger implements migrate.Logger interface
type migrateLogger struct {
	logger *slog.Logger
}

func (l *migrateLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *migrateLogger) Verbose() bool {
	return true
}

// printUsage prints the CLI usage information
func printUsage() {
	usage := `
MediSync Database Migration Tool

Usage:
  migrate <command> [options] [arguments]

Commands:
  up [version]     Run migrations up to the latest or specified version
  down [steps]     Roll back migrations (default: 1 step)
  down-all         Roll back all migrations (WARNING: destroys data)
  version          Show current migration version
  status           Show detailed migration status
  force <version>  Force set migration version (use -1 for no migrations)
  create <name>    Create new migration files
  drop             Drop all database objects (DANGEROUS)
  help             Show this help message

Options:
  --database-url   PostgreSQL connection string (default: $DATABASE_URL)
  --path           Path to migrations directory (default: migrations)
  --verbose        Enable verbose logging

Environment Variables:
  DATABASE_URL     PostgreSQL connection string
                   Format: postgres://user:password@host:port/dbname?sslmode=disable
  LOG_LEVEL        Set to "debug" for debug logging

Examples:
  # Run all pending migrations
  DATABASE_URL="postgres://user:pass@localhost:5432/medisync?sslmode=disable" migrate up

  # Roll back the last migration
  migrate down --database-url="postgres://localhost/medisync"

  # Show migration status
  migrate status

  # Create a new migration
  migrate create add_appointments_table

  # Migrate to a specific version
  migrate up 3

  # Force database to clean state (use after failed migration)
  migrate force 1

Exit Codes:
  0  Success
  1  Error
  2  Usage error
  3  No changes (database already up to date)
  4  Database in dirty state (requires force)
`
	fmt.Println(usage)
}
