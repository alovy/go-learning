package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	"url.com/data/internal/config"
	"url.com/data/internal/crawler"
	"url.com/data/internal/migrate"
)

func main() {

	dbUser := config.GetEnv("DB_USER")
	dbPassword := config.GetEnv("DB_PASS")
	dbName := config.GetEnv("DB_NAME")
	dbHost := config.GetEnvWithDefault("DB_HOST", "db")
	dbPort := config.GetEnvWithDefault("DB_PORT", "5432")

	// Create connection string using environment variables
	dbURL := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", dbUser, dbPassword, dbName, dbHost, dbPort)
	fmt.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// Ensure the database connection is valid
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping the database: %v", err)
	}

	// Run database migrations
	if err := migrate.Do(db); err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	urlsFlag := flag.String("urls", "", "Comma-separated list of URLs to fetch")
	flag.Parse()
	if urlsFlag == nil || *urlsFlag == "" {
		fmt.Println("Please provide URLs with the --urls flag.")
		return
	}

	// Split the URLs by commas
	urls := strings.Split(*urlsFlag, ",")

	var successCount, failureCount int

	var results = make(chan error, len(urls))

	for _, url := range urls {
		go func(url string) {
			err := crawler.Do(url, db)
			results <- err
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		err := <-results
		if err != nil {
			failureCount++
			log.Printf("Error URL: %v\n", err)
		} else {
			successCount++
		}
	}
	fmt.Printf("Success count = %d, Failurecount = %d", successCount, failureCount)

}
