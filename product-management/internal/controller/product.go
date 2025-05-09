package controller

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"

	"product-api/internal/model"

	"github.com/go-chi/chi/v5"
)

type product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) *product {
	return &product{
		db: db,
	}
}

func (p *product) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/", p.Create)
	r.Get("/", p.Get)
	r.Get("/{id}", p.GetProduct)
	r.Put("/{id}", p.Update)
	r.Delete("/{id}", p.Delete)

	return r
}

// Create function inserts a new product into the database.
func (p *product) Create(w http.ResponseWriter, r *http.Request) {
	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := model.CreateProduct(p.db, &product); err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set(contentTypeHeader, applicationJSON)
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(product)
}

func (p *product) Get(w http.ResponseWriter, r *http.Request) {
	page, limit := parsePaginationParams(r)
	offset := (page - 1) * limit

	total, err := model.GetTotalProductsCount(p.db)
	if err != nil {
		http.Error(w, "Error counting products", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	products, err := model.FetchProducts(p.db, limit, offset)
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := prepareResponse(products, page, limit, total)
	w.Header().Set(contentTypeHeader, applicationJSON)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println(err)
	}
}

func parsePaginationParams(r *http.Request) (int, int) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	return page, limit
}

// GetProductByID function retrieves a product by its ID from the database.
func (p *product) GetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := model.FetchProductByID(p.db, productID)
	if err != nil {
		http.Error(w, "Error fetching the product", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (p *product) Update(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product model.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	updated, err := model.UpdateProduct(p.db, productID, &product)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !updated {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *product) Delete(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	deleted, err := model.DeleteProduct(p.db, productID)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if !deleted {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

func prepareResponse(products []model.Product, page, limit, total int) map[string]interface{} {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return map[string]interface{}{
		"data": products,
		"meta": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	}
}
