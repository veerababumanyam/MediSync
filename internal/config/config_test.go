package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Set up test environment variables
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)

	// Set minimal required environment variables
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb?sslmode=disable")
	os.Setenv("NATS_URL", "nats://localhost:4222")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}
}

func TestLoadWithDefaults(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)
	clearEnv()

	// Set only minimal required vars
	os.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/testdb")
	os.Setenv("NATS_URL", "nats://localhost:4222")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Check defaults are applied
	if cfg.App.Environment != EnvDevelopment {
		t.Errorf("Expected environment to be development, got %s", cfg.App.Environment)
	}

	if cfg.App.LogLevel != "info" {
		t.Errorf("Expected log level to be 'info', got %s", cfg.App.LogLevel)
	}

	if cfg.App.LogFormat != "json" {
		t.Errorf("Expected log format to be 'json', got %s", cfg.App.LogFormat)
	}

	if cfg.Database.Port != 5432 {
		t.Errorf("Expected default database port 5432, got %d", cfg.Database.Port)
	}

	if cfg.ETL.BatchSize != 1000 {
		t.Errorf("Expected default batch size 1000, got %d", cfg.ETL.BatchSize)
	}

	if cfg.ETL.SyncInterval != 5*time.Minute {
		t.Errorf("Expected default sync interval 5m, got %v", cfg.ETL.SyncInterval)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantError bool
	}{
		{
			name: "valid config with URL",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/test"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: false,
		},
		{
			name: "valid config with host",
			cfg: &Config{
				Database: DatabaseConfig{Host: "localhost"},
				NATS:     NATSConfig{Host: "localhost"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: false,
		},
		{
			name: "missing database config",
			cfg: &Config{
				Database: DatabaseConfig{},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: true,
		},
		{
			name: "missing NATS config",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/test"},
				NATS:     NATSConfig{},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: true,
		},
		{
			name: "invalid sync interval",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/test"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				ETL:      ETLConfig{SyncInterval: 100 * time.Millisecond, BatchSize: 100},
			},
			wantError: true,
		},
		{
			name: "invalid batch size",
			cfg: &Config{
				Database: DatabaseConfig{URL: "postgres://localhost/test"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 0},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateForProduction(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *Config
		wantError bool
	}{
		{
			name: "valid production config",
			cfg: &Config{
				App:      AppConfig{Environment: EnvProduction},
				Database: DatabaseConfig{URL: "postgres://localhost/test", SSLMode: "require"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				Redis:    RedisConfig{Password: "secret"},
				Keycloak: KeycloakConfig{ClientSecret: "secret"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: false,
		},
		{
			name: "non-production environment",
			cfg: &Config{
				App:      AppConfig{Environment: EnvDevelopment},
				Database: DatabaseConfig{URL: "postgres://localhost/test", SSLMode: "require"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				Redis:    RedisConfig{Password: "secret"},
				Keycloak: KeycloakConfig{ClientSecret: "secret"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: true,
		},
		{
			name: "SSL disabled in production",
			cfg: &Config{
				App:      AppConfig{Environment: EnvProduction},
				Database: DatabaseConfig{URL: "postgres://localhost/test", SSLMode: "disable"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				Redis:    RedisConfig{Password: "secret"},
				Keycloak: KeycloakConfig{ClientSecret: "secret"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: true,
		},
		{
			name: "missing Redis password in production",
			cfg: &Config{
				App:      AppConfig{Environment: EnvProduction},
				Database: DatabaseConfig{URL: "postgres://localhost/test", SSLMode: "require"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				Redis:    RedisConfig{},
				Keycloak: KeycloakConfig{ClientSecret: "secret"},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: true,
		},
		{
			name: "missing Keycloak secret in production",
			cfg: &Config{
				App:      AppConfig{Environment: EnvProduction},
				Database: DatabaseConfig{URL: "postgres://localhost/test", SSLMode: "require"},
				NATS:     NATSConfig{URL: "nats://localhost:4222"},
				Redis:    RedisConfig{Password: "secret"},
				Keycloak: KeycloakConfig{},
				ETL:      ETLConfig{SyncInterval: 5 * time.Minute, BatchSize: 100},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.ValidateForProduction()
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateForProduction() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestDatabaseDSN(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		expected string
	}{
		{
			name: "uses URL if set",
			cfg: &Config{
				Database: DatabaseConfig{
					URL:  "postgres://user:pass@host:5432/db",
					Host: "other",
				},
			},
			expected: "postgres://user:pass@host:5432/db",
		},
		{
			name: "builds DSN from components",
			cfg: &Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "medisync",
					Password: "password",
					Name:     "medisync",
					SSLMode:  "disable",
				},
			},
			expected: "postgres://medisync:password@localhost:5432/medisync?sslmode=disable",
		},
		{
			name: "escapes special characters in password",
			cfg: &Config{
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "user",
					Password: "p@ss:word/test",
					Name:     "db",
					SSLMode:  "require",
				},
			},
			expected: "postgres://user:p%40ss%3Aword%2Ftest@localhost:5432/db?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cfg.DatabaseDSN()
			if result != tt.expected {
				t.Errorf("DatabaseDSN() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestRedisDSN(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		expected string
	}{
		{
			name: "uses URL if set",
			cfg: &Config{
				Redis: RedisConfig{
					URL:  "redis://localhost:6379/0",
					Host: "other",
				},
			},
			expected: "redis://localhost:6379/0",
		},
		{
			name: "builds DSN without password",
			cfg: &Config{
				Redis: RedisConfig{
					Host:     "localhost",
					Port:     6379,
					Database: 0,
				},
			},
			expected: "redis://localhost:6379/0",
		},
		{
			name: "builds DSN with password",
			cfg: &Config{
				Redis: RedisConfig{
					Host:     "localhost",
					Port:     6379,
					Password: "secret",
					Database: 1,
				},
			},
			expected: "redis://:secret@localhost:6379/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cfg.RedisDSN()
			if result != tt.expected {
				t.Errorf("RedisDSN() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestTallyURL(t *testing.T) {
	cfg := &Config{
		Tally: TallyConfig{
			Host: "192.168.1.100",
			Port: 9000,
		},
	}

	expected := "http://192.168.1.100:9000"
	result := cfg.TallyURL()

	if result != expected {
		t.Errorf("TallyURL() = %q, want %q", result, expected)
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env      Environment
		expected bool
	}{
		{EnvProduction, true},
		{EnvStaging, false},
		{EnvDevelopment, false},
	}

	for _, tt := range tests {
		cfg := &Config{App: AppConfig{Environment: tt.env}}
		if result := cfg.IsProduction(); result != tt.expected {
			t.Errorf("IsProduction() for %s = %v, want %v", tt.env, result, tt.expected)
		}
	}
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		env      Environment
		expected bool
	}{
		{EnvDevelopment, true},
		{EnvStaging, false},
		{EnvProduction, false},
	}

	for _, tt := range tests {
		cfg := &Config{App: AppConfig{Environment: tt.env}}
		if result := cfg.IsDevelopment(); result != tt.expected {
			t.Errorf("IsDevelopment() for %s = %v, want %v", tt.env, result, tt.expected)
		}
	}
}

func TestParseEnvironment(t *testing.T) {
	tests := []struct {
		input    string
		expected Environment
	}{
		{"development", EnvDevelopment},
		{"dev", EnvDevelopment},
		{"staging", EnvStaging},
		{"stage", EnvStaging},
		{"production", EnvProduction},
		{"prod", EnvProduction},
		{"PRODUCTION", EnvProduction},
		{"Production", EnvProduction},
		{"unknown", EnvDevelopment},
		{"", EnvDevelopment},
	}

	for _, tt := range tests {
		result := parseEnvironment(tt.input)
		if result != tt.expected {
			t.Errorf("parseEnvironment(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestGetEnv(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)

	os.Setenv("TEST_VAR", "test_value")
	os.Unsetenv("UNSET_VAR")

	tests := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{"TEST_VAR", "default", "test_value"},
		{"UNSET_VAR", "default", "default"},
	}

	for _, tt := range tests {
		result := getEnv(tt.key, tt.defaultValue)
		if result != tt.expected {
			t.Errorf("getEnv(%q, %q) = %q, want %q", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestGetEnvInt(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)

	os.Setenv("INT_VAR", "42")
	os.Setenv("INVALID_INT", "not_a_number")
	os.Unsetenv("UNSET_INT")

	tests := []struct {
		key          string
		defaultValue int
		expected     int
	}{
		{"INT_VAR", 0, 42},
		{"INVALID_INT", 10, 10},
		{"UNSET_INT", 100, 100},
	}

	for _, tt := range tests {
		result := getEnvInt(tt.key, tt.defaultValue)
		if result != tt.expected {
			t.Errorf("getEnvInt(%q, %d) = %d, want %d", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestGetEnvBool(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)

	os.Setenv("BOOL_TRUE", "true")
	os.Setenv("BOOL_FALSE", "false")
	os.Setenv("BOOL_1", "1")
	os.Setenv("BOOL_INVALID", "invalid")
	os.Unsetenv("UNSET_BOOL")

	tests := []struct {
		key          string
		defaultValue bool
		expected     bool
	}{
		{"BOOL_TRUE", false, true},
		{"BOOL_FALSE", true, false},
		{"BOOL_1", false, true},
		{"BOOL_INVALID", true, true},
		{"UNSET_BOOL", false, false},
	}

	for _, tt := range tests {
		result := getEnvBool(tt.key, tt.defaultValue)
		if result != tt.expected {
			t.Errorf("getEnvBool(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestGetEnvDuration(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)

	os.Setenv("DUR_5M", "5m")
	os.Setenv("DUR_1H30M", "1h30m")
	os.Setenv("DUR_300S", "300s")
	os.Setenv("DUR_INVALID", "invalid")
	os.Unsetenv("UNSET_DUR")

	tests := []struct {
		key          string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{"DUR_5M", 0, 5 * time.Minute},
		{"DUR_1H30M", 0, 90 * time.Minute},
		{"DUR_300S", 0, 300 * time.Second},
		{"DUR_INVALID", time.Hour, time.Hour},
		{"UNSET_DUR", 10 * time.Second, 10 * time.Second},
	}

	for _, tt := range tests {
		result := getEnvDuration(tt.key, tt.defaultValue)
		if result != tt.expected {
			t.Errorf("getEnvDuration(%q, %v) = %v, want %v", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestMustLoad_Panics(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)
	clearEnv()

	// Set invalid batch size to trigger validation error
	// (batch size must be at least 1)
	os.Setenv("ETL_BATCH_SIZE", "0")

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLoad() did not panic on invalid config")
		}
	}()

	MustLoad()
}

func TestMustLoad_Success(t *testing.T) {
	originalEnv := os.Environ()
	defer restoreEnv(originalEnv)
	clearEnv()

	// With defaults, Load() should succeed
	// (defaults provide localhost for DB and NATS)

	cfg := MustLoad()
	if cfg == nil {
		t.Error("MustLoad() returned nil config")
	}
}

// Helper functions for tests

func clearEnv() {
	envVars := []string{
		"APP_ENV", "LOG_LEVEL", "LOG_FORMAT", "TIMEZONE",
		"DATABASE_URL", "POSTGRES_HOST", "POSTGRES_PORT", "POSTGRES_USER",
		"POSTGRES_PASSWORD", "POSTGRES_DB", "POSTGRES_SSLMODE",
		"NATS_URL", "NATS_HOST", "NATS_PORT",
		"REDIS_URL", "REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD",
		"KEYCLOAK_URL", "KEYCLOAK_REALM", "KEYCLOAK_CLIENT_ID", "KEYCLOAK_CLIENT_SECRET",
		"TALLY_HOST", "TALLY_PORT", "TALLY_COMPANY",
		"HIMS_API_URL", "HIMS_API_KEY",
		"ETL_METRICS_PORT", "ETL_SYNC_INTERVAL", "ETL_BATCH_SIZE",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}
}

func restoreEnv(originalEnv []string) {
	os.Clearenv()
	for _, e := range originalEnv {
		pair := splitEnvPair(e)
		if len(pair) == 2 {
			os.Setenv(pair[0], pair[1])
		}
	}
}

func splitEnvPair(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return []string{env}
}
