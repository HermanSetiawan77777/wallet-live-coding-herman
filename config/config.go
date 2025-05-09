package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	Port        int
	Environment string
	DB          DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func LoadConfig() (*AppConfig, error) {
	port, err := strconv.Atoi(getEnvOrDefault("APP_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid APP_PORT: %v", err)
	}

	dbConfig, err := LoadDBConfig()
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Port:        port,
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		DB:          dbConfig,
	}, nil
}

func LoadDBConfig() (DBConfig, error) {
	port, err := strconv.Atoi(getEnvOrDefault("DB_PORT", "1433"))
	if err != nil {
		return DBConfig{}, fmt.Errorf("invalid DB_PORT: %v", err)
	}

	return DBConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnvOrDefault("DB_USER", "sa"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   getEnvOrDefault("DB_NAME", "master"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

