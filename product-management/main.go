package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"product-api/internal/database"
	"product-api/internal/middleware"
	"product-api/internal/migrate"
)

func main() {
	db, err := database.InitializeDatabase()
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := migrate.Do(db); err != nil {
		log.Fatalf("failed to setup database: %v", err)
	}

	r := middleware.SetupRouter(db)

	// Start the server
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
