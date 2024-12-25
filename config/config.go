package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig   `envconfig:"SERVER"`
	Database DatabaseConfig `envconfig:"DB"`
	Storage  StorageConfig  `envconfig:"STORAGE"`
}

type ServerConfig struct {
	Port int    `envconfig:"PORT" default:"8080"`
	Mode string `envconfig:"MODE" default:"dev"` // `release`, `dev`, or `test`
}

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	UseCloudSQL bool           `envconfig:"USE_CLOUD_SQL"`
	Master      DatabaseDetail `envconfig:"MASTER"`
	Replica     DatabaseDetail `envconfig:"REPLICA"`
}

type DatabaseDetail struct {
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	DBName   string `envconfig:"NAME"`
}

type StorageConfig struct {
	RootFolder string `envconfig:"ROOT_FOLDER"`
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() (*Config, error) {
	// load .env file if exist
	godotenv.Load()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from environment variables: %v", err)
	}
	return &cfg, nil
}
