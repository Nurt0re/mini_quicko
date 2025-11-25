package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Kaspi    KaspiConfig    `yaml:"kaspi"`
	App      AppConfig      `yaml:"app"`
}

type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"-"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type KaspiConfig struct {
	DefaultCityID string `yaml:"default_city_id"`
}

type AppConfig struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func Load(configPath string) (*Config, error) {

	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using defaults")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD not set in .env file")
	}

	return &cfg, nil
}

func (c *ServerConfig) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
