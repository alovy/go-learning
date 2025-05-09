package config

import (
	"log"
	"os"
)

func GetEnv(key, defaultValue string) string {
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
