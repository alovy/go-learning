package database

import (
	"database/sql"
	"fmt"
	"log"
	"product-api/internal/config"
)

func InitializeDatabase() (*sql.DB, error) {
	dbUser := config.GetEnv("DB_USER", "")
	dbPassword := config.GetEnv("DB_PASS", "")
	dbName := config.GetEnv("DB_NAME", "")
	dbHost := config.GetEnv("DB_HOST", "db")
	dbPort := config.GetEnv("DB_PORT", "5432")

	// Create connection string using environment variables
	dbURL := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
	fmt.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
		return nil, err
	}

	// Ensure the database connection is valid
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping the database: %v", err)
		return nil, err
	}
	return db, nil
}
