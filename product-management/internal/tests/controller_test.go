package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"product-api/internal/controller"
	"product-api/internal/jwt"
	"product-api/internal/middleware"
	"product-api/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	// Set fake env vars required by your app
	os.Setenv("DB_PASS", "dummy")
	os.Setenv("JWT_EXPIRATION_TIME", "300s")
	os.Setenv("JWT_SECRET_KEY", "testsecret")
	os.Setenv("JWT_USERNAME", "testuser")

	// Run the tests
	code := m.Run()

	// Clean up or unset env vars if needed
	os.Exit(code)
}

func TestCreateProductHandler(t *testing.T) {
	mockService := new(MockProductService)

	mockService.On("CreateProduct", mock.AnythingOfType("*model.Product")).Return(nil)

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/products/", strings.NewReader(`{
		"name": "Test Product",
		"category": "Test",
		"description": "Test Desc",
		"price": 12.34
	}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Mount("/", controller.NewProduct(mockService).Router())
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"Test Product"`)

	mockService.AssertExpectations(t)
}

func TestGetProductsHandler(t *testing.T) {

	mockService := new(MockProductService)

	products := []model.Product{
		{
			ProductID:   "1",
			Name:        "Test Product",
			Category:    "Test",
			Description: "Test Desc",
			Price:       12.34,
		},
	}
	totalCount := 1

	mockService.On("GetAllProducts", 10, 0).Return(products, totalCount, nil)

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Mount("/", controller.NewProduct(mockService).Router())
	})

	req := httptest.NewRequest(http.MethodGet, "/products/?page=1&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"Test Product"`)

	mockService.AssertExpectations(t)
}

func TestGetProductByIDHandler(t *testing.T) {
	mockService := new(MockProductService)

	product := &model.Product{
		ProductID:   "1",
		Name:        "Test Product",
		Category:    "Test",
		Description: "Test Desc",
		Price:       12.34,
	}

	mockService.On("GetProductByID", 1).Return(product, nil)

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Mount("/", controller.NewProduct(mockService).Router())
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"name":"Test Product"`)

	mockService.AssertExpectations(t)
}

func TestUpdateProductHandler(t *testing.T) {
	mockService := new(MockProductService)

	mockService.On("UpdateProduct", 1, mock.AnythingOfType("*model.Product")).Return(true, nil)

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(`{
		"name": "Updated Product",
		"category": "Updated Category",
		"description": "Updated Desc",
		"price": 56.78
	}`))
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Mount("/", controller.NewProduct(mockService).Router())
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockService.AssertExpectations(t)
}

func TestDeleteProductHandler(t *testing.T) {
	mockService := new(MockProductService)

	mockService.On("DeleteProduct", 1).Return(true, nil)

	token, err := jwt.GenerateToken()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()

	router := chi.NewRouter()
	router.Route("/products", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Mount("/", controller.NewProduct(mockService).Router())
	})

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	mockService.AssertExpectations(t)
}
