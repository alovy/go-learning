package config

import (
	"log"
	"os"
)

type DB struct {
	Name     string
	Host     string
	Port     string
	Username string
	Password string
}

type TOKEN struct {
	Secret string
	Expiry string
	User   string
}

type Config struct {
	DB
	TOKEN
}

func LoadConfig() Config {
	db := DB{
		Username: getEnv("DB_USER", ""),
		Password: getEnv("DB_PASS", ""),
		Name:     getEnv("DB_NAME", ""),
		Host:     getEnv("DB_HOST", "db"),
		Port:     getEnv("DB_PORT", "5432"),
	}
	token := TOKEN{
		Secret: getEnv("JWT_SECRET_KEY", ""),
		Expiry: getEnv("JWT_EXPIRATION_TIME", "24h"),
		User:   getEnv("JWT_USERNAME", ""),
	}
	return Config{
		DB:    db,
		TOKEN: token,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		if defaultValue == "" {
			log.Fatalf("environment variable %s is not set", key)
		} else {
			return defaultValue
		}
	}
	return value
}
