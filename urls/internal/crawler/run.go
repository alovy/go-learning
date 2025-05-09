package crawler

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbConnStr = "user=postgres dbname=mydb password=password host=db port=5432 sslmode=disable"
)

var db *sql.DB

func Do(url string, db *sql.DB) error {
	data, err := FetchURL(url, db)
	if err != nil {
		return err
	}
	err = SaveURL(url, data, db)
	if err != nil {
		return err
	}
	return nil
}

// FetchURL fetches a URL, processes the response, and stores it in the database
func FetchURL(url string, db *sql.DB) ([]byte, error) {

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("failed to fetch URL %s: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body from %s: %v", url, err)
		return nil, err
	}
	return body, nil
}

// Insert the URL and response into the database
func SaveURL(url string, data []byte, db *sql.DB) error {
	const insertURLResponseQuery = `INSERT INTO url_responses (url, response) VALUES ($1, $2)`
	_, err := db.Exec(insertURLResponseQuery, url, string(data))

	if err != nil {
		log.Printf("failed to insert URL %s into database: %v", url, err)
		return err
	}
	return nil
}
