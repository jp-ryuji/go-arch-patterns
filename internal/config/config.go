package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	// Database configuration
	DBHost            string        `mapstructure:"DB_HOST"`
	DBPort            int           `mapstructure:"DB_PORT"`
	DBPortExternal    int           `mapstructure:"DB_PORT_EXTERNAL"`
	DBUser            string        `mapstructure:"DB_USER"`
	DBPassword        string        `mapstructure:"DB_PASSWORD"`
	DBName            string        `mapstructure:"DB_NAME"`
	DBSSLMode         string        `mapstructure:"DB_SSLMODE"`
	DBMaxOpenConns    int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	DBMaxIdleConns    int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBConnMaxLifetime time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`

	// Redis configuration
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort int    `mapstructure:"REDIS_PORT"`
	RedisURL  string `mapstructure:"REDIS_URL"`

	// Server configuration
	GRPCPort int `mapstructure:"GRPC_PORT"`
	HTTPPort int `mapstructure:"HTTP_PORT"`

	// OpenSearch configuration
	OpenSearchPortExternal int `mapstructure:"OPENSEARCH_PORT_EXTERNAL"`
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	// Initialize Viper
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Set default values
	setDefaults()

	// Bind environment variables
	bindEnv()

	// Create config struct
	var config Config

	// Unmarshal config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration options
func setDefaults() {
	// Database defaults
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_PORT_EXTERNAL", 5434)
	viper.SetDefault("DB_USER", "user")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "mydb")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 25)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute)

	// Redis defaults
	viper.SetDefault("REDIS_HOST", "redis")
	viper.SetDefault("REDIS_PORT", 6379)
	viper.SetDefault("REDIS_URL", "redis://redis:6379")

	// Server defaults
	viper.SetDefault("GRPC_PORT", 50051)
	viper.SetDefault("HTTP_PORT", 8081)

	// OpenSearch defaults
	viper.SetDefault("OPENSEARCH_PORT_EXTERNAL", 9201)
}

// bindEnv binds environment variables to Viper keys
func bindEnv() {
	// Database
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")
	_ = viper.BindEnv("DB_PORT_EXTERNAL")
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("DB_SSLMODE")
	_ = viper.BindEnv("DB_MAX_OPEN_CONNS")
	_ = viper.BindEnv("DB_MAX_IDLE_CONNS")
	_ = viper.BindEnv("DB_CONN_MAX_LIFETIME")

	// Redis
	_ = viper.BindEnv("REDIS_HOST")
	_ = viper.BindEnv("REDIS_PORT")
	_ = viper.BindEnv("REDIS_URL")

	// Server
	_ = viper.BindEnv("GRPC_PORT")
	_ = viper.BindEnv("HTTP_PORT")

	// OpenSearch
	_ = viper.BindEnv("OPENSEARCH_PORT_EXTERNAL")
}

// DatabaseURL returns the database connection string
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPortExternal, c.DBName, c.DBSSLMode)
}
