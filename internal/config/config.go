// Package config provides environment configuration loading for MediSync services.
//
// Configuration is loaded from environment variables with sensible defaults for
// development. All services (PostgreSQL, NATS, Redis, Keycloak, Tally, HIMS)
// are configured through this package.
//
// Usage:
//
//	cfg, err := config.Load()
//	if err != nil {
//	    log.Fatal("Failed to load configuration:", err)
//	}
//
// Environment variables can be set directly or loaded from a .env file.
// See .env.example in the project root for all available configuration options.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Environment represents the application environment.
type Environment string

const (
	// EnvDevelopment indicates a development environment.
	EnvDevelopment Environment = "development"
	// EnvStaging indicates a staging environment.
	EnvStaging Environment = "staging"
	// EnvProduction indicates a production environment.
	EnvProduction Environment = "production"
)

// Config holds all application configuration.
type Config struct {
	// Application settings
	App AppConfig

	// Database configuration
	Database DatabaseConfig

	// NATS messaging configuration
	NATS NATSConfig

	// Redis cache configuration
	Redis RedisConfig

	// Keycloak authentication configuration
	Keycloak KeycloakConfig

	// Tally ERP connection configuration
	Tally TallyConfig

	// HIMS API connection configuration
	HIMS HIMSConfig

	// ETL service configuration
	ETL ETLConfig

	// Observability configuration
	Observability ObservabilityConfig

	// LLM configuration
	LLM LLMConfig

	// OPA configuration
	OPA OPAConfig

	// Server configuration
	Server ServerConfig
}

// AppConfig holds general application settings.
type AppConfig struct {
	// Environment is the application environment (development, staging, production).
	Environment Environment

	// LogLevel is the logging level (debug, info, warn, error).
	LogLevel string

	// LogFormat is the log output format (json, text).
	LogFormat string

	// Timezone is the application timezone.
	Timezone string
}

// DatabaseConfig holds PostgreSQL connection settings.
type DatabaseConfig struct {
	// URL is the full PostgreSQL connection string.
	URL string

	// Host is the database server hostname.
	Host string

	// Port is the database server port.
	Port int

	// User is the database username.
	User string

	// Password is the database password.
	Password string

	// Name is the database name.
	Name string

	// SSLMode is the SSL connection mode (disable, require, verify-ca, verify-full).
	SSLMode string

	// MaxOpenConns is the maximum number of open connections.
	MaxOpenConns int

	// MaxIdleConns is the maximum number of idle connections.
	MaxIdleConns int

	// ConnMaxLifetime is the maximum connection lifetime.
	ConnMaxLifetime time.Duration

	// ConnMaxIdleTime is the maximum connection idle time.
	ConnMaxIdleTime time.Duration
}

// NATSConfig holds NATS messaging settings.
type NATSConfig struct {
	// URL is the NATS server URL.
	URL string

	// Host is the NATS server hostname.
	Host string

	// Port is the NATS client port.
	Port int

	// MonitorPort is the NATS HTTP monitoring port.
	MonitorPort int

	// ClusterPort is the NATS cluster routing port.
	ClusterPort int

	// MaxReconnects is the maximum number of reconnection attempts.
	MaxReconnects int

	// ReconnectWait is the wait duration between reconnection attempts.
	ReconnectWait time.Duration

	// JetStreamEnabled indicates if JetStream is enabled.
	JetStreamEnabled bool
}

// RedisConfig holds Redis cache settings.
type RedisConfig struct {
	// URL is the full Redis connection URL.
	URL string

	// Host is the Redis server hostname.
	Host string

	// Port is the Redis server port.
	Port int

	// Password is the Redis password (optional).
	Password string

	// Database is the Redis database number.
	Database int

	// MaxRetries is the maximum number of retries.
	MaxRetries int

	// PoolSize is the connection pool size.
	PoolSize int

	// MinIdleConns is the minimum number of idle connections.
	MinIdleConns int

	// DialTimeout is the connection timeout.
	DialTimeout time.Duration

	// ReadTimeout is the read operation timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the write operation timeout.
	WriteTimeout time.Duration
}

// KeycloakConfig holds Keycloak authentication settings.
type KeycloakConfig struct {
	// URL is the Keycloak server base URL.
	URL string

	// Realm is the Keycloak realm name.
	Realm string

	// ClientID is the OAuth2 client ID.
	ClientID string

	// ClientSecret is the OAuth2 client secret (optional for public clients).
	ClientSecret string

	// AdminUser is the Keycloak admin username.
	AdminUser string

	// AdminPassword is the Keycloak admin password.
	AdminPassword string
}

// TallyConfig holds Tally ERP connection settings.
type TallyConfig struct {
	// Host is the Tally server hostname.
	Host string

	// Port is the Tally server port.
	Port int

	// Company is the Tally company name.
	Company string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int

	// RetryDelay is the initial delay between retries (exponential backoff).
	RetryDelay time.Duration
}

// HIMSConfig holds HIMS API connection settings.
type HIMSConfig struct {
	// URL is the HIMS API base URL.
	URL string

	// APIKey is the HIMS API authentication key.
	APIKey string

	// Timeout is the HTTP request timeout.
	Timeout time.Duration

	// MaxRetries is the maximum number of retry attempts.
	MaxRetries int

	// RetryDelay is the initial delay between retries (exponential backoff).
	RetryDelay time.Duration

	// RateLimitPerSecond is the maximum requests per second.
	RateLimitPerSecond int
}

// ETLConfig holds ETL service configuration.
type ETLConfig struct {
	// MetricsPort is the Prometheus metrics HTTP port.
	MetricsPort int

	// SyncInterval is the interval between sync runs.
	SyncInterval time.Duration

	// BatchSize is the number of records to process per batch.
	BatchSize int

	// WorkerCount is the number of concurrent workers.
	WorkerCount int

	// QuarantineEnabled enables quarantine for failed records.
	QuarantineEnabled bool

	// RetryFailedRecords enables automatic retry of quarantined records.
	RetryFailedRecords bool

	// MaxRetries is the maximum retry attempts for failed records.
	MaxRetries int
}

// ObservabilityConfig holds monitoring and logging settings.
type ObservabilityConfig struct {
	// PrometheusPort is the Prometheus server port.
	PrometheusPort int

	// GrafanaPort is the Grafana server port.
	GrafanaPort int

	// TracingEnabled enables distributed tracing.
	TracingEnabled bool

	// TracingEndpoint is the tracing collector endpoint.
	TracingEndpoint string

	// MetricsEnabled enables Prometheus metrics.
	MetricsEnabled bool
}

// LLMConfig holds LLM provider settings.
type LLMConfig struct {
	// Provider is the LLM provider (openai, ollama, etc.)
	Provider string

	// Model is the model to use.
	Model string

	// APIKey is the API key for the provider.
	APIKey string

	// BaseURL is the base URL for the provider API.
	BaseURL string

	// MaxTokens is the maximum tokens for responses.
	MaxTokens int

	// Temperature is the sampling temperature.
	Temperature float64
}

// OPAConfig holds OPA policy engine settings.
type OPAConfig struct {
	// URL is the OPA server URL.
	URL string

	// Timeout is the request timeout.
	Timeout time.Duration

	// DecisionPath is the path for decision API.
	DecisionPath string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	// Port is the server port.
	Port int

	// Host is the server host.
	Host string

	// ReadTimeout is the read timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the write timeout.
	WriteTimeout time.Duration

	// ShutdownTimeout is the graceful shutdown timeout.
	ShutdownTimeout time.Duration
}

// Load reads configuration from environment variables and returns a Config struct.
// It applies sensible defaults for development and validates required fields.
func Load() (*Config, error) {
	cfg := &Config{}

	// Load all configuration sections
	cfg.App = loadAppConfig()
	cfg.Database = loadDatabaseConfig()
	cfg.NATS = loadNATSConfig()
	cfg.Redis = loadRedisConfig()
	cfg.Keycloak = loadKeycloakConfig()
	cfg.Tally = loadTallyConfig()
	cfg.HIMS = loadHIMSConfig()
	cfg.ETL = loadETLConfig()
	cfg.Observability = loadObservabilityConfig()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// MustLoad loads configuration and panics on error.
// Use this for application startup where configuration is required.
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load configuration: %v", err))
	}
	return cfg
}

// Validate checks that all required configuration values are present and valid.
func (c *Config) Validate() error {
	var errs []error

	// Validate database configuration
	if c.Database.URL == "" && c.Database.Host == "" {
		errs = append(errs, errors.New("database: either DATABASE_URL or POSTGRES_HOST must be set"))
	}

	// Validate NATS configuration
	if c.NATS.URL == "" && c.NATS.Host == "" {
		errs = append(errs, errors.New("nats: either NATS_URL or NATS_HOST must be set"))
	}

	// Validate ETL configuration
	if c.ETL.SyncInterval < time.Second {
		errs = append(errs, errors.New("etl: sync interval must be at least 1 second"))
	}

	if c.ETL.BatchSize < 1 {
		errs = append(errs, errors.New("etl: batch size must be at least 1"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// ValidateForProduction performs stricter validation for production environments.
func (c *Config) ValidateForProduction() error {
	if err := c.Validate(); err != nil {
		return err
	}

	var errs []error

	// Production-specific validations
	if c.App.Environment != EnvProduction {
		errs = append(errs, errors.New("app: environment must be 'production' for production deployment"))
	}

	if c.Database.SSLMode == "disable" {
		errs = append(errs, errors.New("database: SSL must be enabled in production"))
	}

	if c.Redis.Password == "" {
		errs = append(errs, errors.New("redis: password must be set in production"))
	}

	if c.Keycloak.ClientSecret == "" {
		errs = append(errs, errors.New("keycloak: client secret must be set in production"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// IsProduction returns true if running in production environment.
func (c *Config) IsProduction() bool {
	return c.App.Environment == EnvProduction
}

// IsDevelopment returns true if running in development environment.
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == EnvDevelopment
}

// DatabaseDSN returns the database connection string.
// If DATABASE_URL is set, it returns that. Otherwise, it constructs the DSN from components.
func (c *Config) DatabaseDSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}

	// Build DSN from components
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(c.Database.User),
		url.QueryEscape(c.Database.Password),
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// RedisDSN returns the Redis connection string.
func (c *Config) RedisDSN() string {
	if c.Redis.URL != "" {
		return c.Redis.URL
	}

	// Build DSN from components
	if c.Redis.Password != "" {
		return fmt.Sprintf("redis://:%s@%s:%d/%d",
			url.QueryEscape(c.Redis.Password),
			c.Redis.Host,
			c.Redis.Port,
			c.Redis.Database,
		)
	}

	return fmt.Sprintf("redis://%s:%d/%d",
		c.Redis.Host,
		c.Redis.Port,
		c.Redis.Database,
	)
}

// TallyURL returns the Tally server URL.
func (c *Config) TallyURL() string {
	return fmt.Sprintf("http://%s:%d", c.Tally.Host, c.Tally.Port)
}

// LogConfig logs the current configuration (with sensitive values masked).
func (c *Config) LogConfig(logger *slog.Logger) {
	logger.Info("Configuration loaded",
		slog.Group("app",
			slog.String("environment", string(c.App.Environment)),
			slog.String("log_level", c.App.LogLevel),
			slog.String("log_format", c.App.LogFormat),
			slog.String("timezone", c.App.Timezone),
		),
		slog.Group("database",
			slog.String("host", c.Database.Host),
			slog.Int("port", c.Database.Port),
			slog.String("name", c.Database.Name),
			slog.String("ssl_mode", c.Database.SSLMode),
			slog.Int("max_open_conns", c.Database.MaxOpenConns),
		),
		slog.Group("nats",
			slog.String("host", c.NATS.Host),
			slog.Int("port", c.NATS.Port),
			slog.Bool("jetstream", c.NATS.JetStreamEnabled),
		),
		slog.Group("redis",
			slog.String("host", c.Redis.Host),
			slog.Int("port", c.Redis.Port),
			slog.Int("database", c.Redis.Database),
		),
		slog.Group("keycloak",
			slog.String("url", c.Keycloak.URL),
			slog.String("realm", c.Keycloak.Realm),
			slog.String("client_id", c.Keycloak.ClientID),
		),
		slog.Group("tally",
			slog.String("host", c.Tally.Host),
			slog.Int("port", c.Tally.Port),
			slog.String("company", c.Tally.Company),
		),
		slog.Group("hims",
			slog.String("url", c.HIMS.URL),
			slog.Bool("api_key_set", c.HIMS.APIKey != ""),
		),
		slog.Group("etl",
			slog.Int("metrics_port", c.ETL.MetricsPort),
			slog.Duration("sync_interval", c.ETL.SyncInterval),
			slog.Int("batch_size", c.ETL.BatchSize),
			slog.Int("worker_count", c.ETL.WorkerCount),
		),
	)
}

// loadAppConfig loads application settings from environment variables.
func loadAppConfig() AppConfig {
	env := getEnv("APP_ENV", "development")

	return AppConfig{
		Environment: parseEnvironment(env),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		LogFormat:   getEnv("LOG_FORMAT", "json"),
		Timezone:    getEnv("TIMEZONE", "UTC"),
	}
}

// loadDatabaseConfig loads PostgreSQL settings from environment variables.
func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		URL:             getEnv("DATABASE_URL", ""),
		Host:            getEnv("POSTGRES_HOST", "localhost"),
		Port:            getEnvInt("POSTGRES_PORT", 5432),
		User:            getEnv("POSTGRES_USER", "medisync"),
		Password:        getEnv("POSTGRES_PASSWORD", "medisync_dev_password"),
		Name:            getEnv("POSTGRES_DB", "medisync"),
		SSLMode:         getEnv("POSTGRES_SSLMODE", "disable"),
		MaxOpenConns:    getEnvInt("POSTGRES_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("POSTGRES_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: getEnvDuration("POSTGRES_CONN_MAX_LIFETIME", 5*time.Minute),
		ConnMaxIdleTime: getEnvDuration("POSTGRES_CONN_MAX_IDLE_TIME", 1*time.Minute),
	}
}

// loadNATSConfig loads NATS settings from environment variables.
func loadNATSConfig() NATSConfig {
	return NATSConfig{
		URL:              getEnv("NATS_URL", ""),
		Host:             getEnv("NATS_HOST", "localhost"),
		Port:             getEnvInt("NATS_PORT", 4222),
		MonitorPort:      getEnvInt("NATS_MONITOR_PORT", 8222),
		ClusterPort:      getEnvInt("NATS_CLUSTER_PORT", 6222),
		MaxReconnects:    getEnvInt("NATS_MAX_RECONNECTS", 10),
		ReconnectWait:    getEnvDuration("NATS_RECONNECT_WAIT", 2*time.Second),
		JetStreamEnabled: getEnvBool("NATS_JETSTREAM_ENABLED", true),
	}
}

// loadRedisConfig loads Redis settings from environment variables.
func loadRedisConfig() RedisConfig {
	return RedisConfig{
		URL:          getEnv("REDIS_URL", ""),
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         getEnvInt("REDIS_PORT", 6379),
		Password:     getEnv("REDIS_PASSWORD", ""),
		Database:     getEnvInt("REDIS_DB", 0),
		MaxRetries:   getEnvInt("REDIS_MAX_RETRIES", 3),
		PoolSize:     getEnvInt("REDIS_POOL_SIZE", 10),
		MinIdleConns: getEnvInt("REDIS_MIN_IDLE_CONNS", 2),
		DialTimeout:  getEnvDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
		ReadTimeout:  getEnvDuration("REDIS_READ_TIMEOUT", 3*time.Second),
		WriteTimeout: getEnvDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
	}
}

// loadKeycloakConfig loads Keycloak settings from environment variables.
func loadKeycloakConfig() KeycloakConfig {
	return KeycloakConfig{
		URL:           getEnv("KEYCLOAK_URL", "http://localhost:8180"),
		Realm:         getEnv("KEYCLOAK_REALM", "medisync"),
		ClientID:      getEnv("KEYCLOAK_CLIENT_ID", "medisync-app"),
		ClientSecret:  getEnv("KEYCLOAK_CLIENT_SECRET", ""),
		AdminUser:     getEnv("KEYCLOAK_ADMIN", "admin"),
		AdminPassword: getEnv("KEYCLOAK_ADMIN_PASSWORD", ""),
	}
}

// loadTallyConfig loads Tally ERP settings from environment variables.
func loadTallyConfig() TallyConfig {
	return TallyConfig{
		Host:       getEnv("TALLY_HOST", "localhost"),
		Port:       getEnvInt("TALLY_PORT", 9000),
		Company:    getEnv("TALLY_COMPANY", ""),
		Timeout:    getEnvDuration("TALLY_TIMEOUT", 30*time.Second),
		MaxRetries: getEnvInt("TALLY_MAX_RETRIES", 3),
		RetryDelay: getEnvDuration("TALLY_RETRY_DELAY", 1*time.Second),
	}
}

// loadHIMSConfig loads HIMS API settings from environment variables.
func loadHIMSConfig() HIMSConfig {
	return HIMSConfig{
		URL:                getEnv("HIMS_API_URL", "http://localhost:8082/api"),
		APIKey:             getEnv("HIMS_API_KEY", ""),
		Timeout:            getEnvDuration("HIMS_TIMEOUT", 30*time.Second),
		MaxRetries:         getEnvInt("HIMS_MAX_RETRIES", 3),
		RetryDelay:         getEnvDuration("HIMS_RETRY_DELAY", 1*time.Second),
		RateLimitPerSecond: getEnvInt("HIMS_RATE_LIMIT", 10),
	}
}

// loadETLConfig loads ETL service settings from environment variables.
func loadETLConfig() ETLConfig {
	return ETLConfig{
		MetricsPort:        getEnvInt("ETL_METRICS_PORT", 9100),
		SyncInterval:       getEnvDuration("ETL_SYNC_INTERVAL", 5*time.Minute),
		BatchSize:          getEnvInt("ETL_BATCH_SIZE", 1000),
		WorkerCount:        getEnvInt("ETL_WORKER_COUNT", 4),
		QuarantineEnabled:  getEnvBool("ETL_QUARANTINE_ENABLED", true),
		RetryFailedRecords: getEnvBool("ETL_RETRY_FAILED", true),
		MaxRetries:         getEnvInt("ETL_MAX_RETRIES", 3),
	}
}

// loadObservabilityConfig loads monitoring settings from environment variables.
func loadObservabilityConfig() ObservabilityConfig {
	return ObservabilityConfig{
		PrometheusPort:  getEnvInt("PROMETHEUS_PORT", 9090),
		GrafanaPort:     getEnvInt("GRAFANA_PORT", 3001),
		TracingEnabled:  getEnvBool("TRACING_ENABLED", false),
		TracingEndpoint: getEnv("TRACING_ENDPOINT", ""),
		MetricsEnabled:  getEnvBool("METRICS_ENABLED", true),
	}
}

// parseEnvironment converts a string to Environment type.
func parseEnvironment(env string) Environment {
	switch strings.ToLower(env) {
	case "production", "prod":
		return EnvProduction
	case "staging", "stage":
		return EnvStaging
	default:
		return EnvDevelopment
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an environment variable as an integer or returns a default value.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool retrieves an environment variable as a boolean or returns a default value.
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// getEnvDuration retrieves an environment variable as a duration or returns a default value.
// Supports Go duration strings (e.g., "5m", "1h30m", "300s").
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
