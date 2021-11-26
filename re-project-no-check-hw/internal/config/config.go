package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type DatabaseConfig interface {
	Uri() string
}

type PostgresConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Dbname   string
}

func (pc PostgresConfig) Uri() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		pc.User, pc.Password, pc.Host, pc.Port, pc.Dbname)
}

type ServerConfig struct {
	Address string
}

type RedisConfig struct {
	Host    string
	Db      int
	Expires time.Duration
}

type Config struct {
	Database DatabaseConfig
	Redis RedisConfig
	Server ServerConfig
}

func NewConfig() *Config {
	return &Config{
		Database: PostgresConfig{
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnvAsInt("POSTGRES_PORT", 5432),
			Dbname:   getEnv("POSTGRES_DBNAME", "postgres"),
		},
		Redis: RedisConfig{
			Host:    getEnv("REDIS_HOST", "localhost:6379"),
			Db:      getEnvAsInt("REDIS_DB", 0),
			Expires: getEnvAsTimeDuration("REDIS_EXPIRES", 0),
		},
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsTimeDuration(key string, defaultVal time.Duration) time.Duration {
	valueStr := getEnv(key, "0")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultVal
}
