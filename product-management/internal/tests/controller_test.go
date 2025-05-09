package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"product-api/internal/controller"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	body := `{
		"name": "Test Product",
		"category": "Test",
		"description": "Test Desc",
		"price": 12.34
	}`

	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs("Test Product", "Test", "Test Desc", 12.34).
		WillReturnRows(sqlmock.NewRows([]string{"product_id"}).AddRow("1"))

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := controller.CreateProduct(db)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"product_id":"1"`)
}

func TestGetProductsHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM products`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products`).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "category", "description", "price"}).
			AddRow("1", "Test Product", "Test", "Test Desc", 12.34))

	req := httptest.NewRequest(http.MethodGet, "/products?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	handler := controller.GetProduct(db)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"product_id":"1"`)
}

func TestGetProductByIDHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products WHERE product_id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "category", "description", "price"}).
			AddRow("1", "Test Product", "Test", "Test Desc", 12.34))

	req := httptest.NewRequest(http.MethodGet, "/products?product_id=1", nil)
	w := httptest.NewRecorder()

	handler := controller.GetProductById(db, 1)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"product_id":"1"`)
}

func TestUpdateProductHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	body := `{
		"name": "Updated Product",
		"category": "Updated",
		"description": "Updated Desc",
		"price": 56.78
	}`

	mock.ExpectExec(`UPDATE products SET name = \$1, category = \$2, description = \$3, price = \$4 WHERE product_id = \$5`).
		WithArgs("Updated Product", "Updated", "Updated Desc", 56.78, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodPut, "/products?product_id=1", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := controller.Update(db, 1)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteProductHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(`DELETE FROM products WHERE product_id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	req := httptest.NewRequest(http.MethodDelete, "/products?product_id=1", nil)
	w := httptest.NewRecorder()

	handler := controller.DeleteProduct(db, 1)
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
