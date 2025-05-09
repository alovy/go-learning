package crawler

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchURL(t *testing.T) {

	tests := []struct {
		name             string
		url              string
		mockHttpResponse string
		expectedResult   error
		mockSaveURLError bool
	}{
		{
			name:             "Successful Fetch",
			url:              "https://jsonplaceholder.typicode.com/todos/1",
			mockHttpResponse: `{"id": 1, "title": "Test Todo"}`,
			expectedResult:   nil,
			mockSaveURLError: false,
		},
		{
			name:             "Failed Fetch - HTTP Error",
			url:              "https://jsonplaceholder.typicode.com/todos/2",
			mockHttpResponse: `{"id": 2, "title": "Test Todo"}`,
			expectedResult:   fmt.Errorf("failed to insert into database"),
			mockSaveURLError: true,
		},
	}

	// Activate HTTP mock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Mock the HTTP response for the given URL
			httpmock.RegisterResponder("GET", testCase.url,
				httpmock.NewStringResponder(200, testCase.mockHttpResponse))

			// Create a mock database connection
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			// Setup mock database interaction for SaveURL
			if testCase.mockSaveURLError {
				// Simulate a database error (failed insert)
				mock.ExpectExec(`INSERT INTO url_responses`).
					WithArgs(testCase.url, testCase.mockHttpResponse).
					WillReturnError(fmt.Errorf("failed to insert into database"))
			} else {
				// Simulate successful insert into the database
				mock.ExpectExec(`INSERT INTO url_responses`).
					WithArgs(testCase.url, testCase.mockHttpResponse).
					WillReturnResult(sqlmock.NewResult(1, 1))
			}

			result := Do(testCase.url, db)

			// Assert the result is as expected
			assert.Equal(t, testCase.expectedResult, result)

			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("There were unmet expectations: %v", err)
			}
		})
	}
}
