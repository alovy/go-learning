package tests

import (
	"database/sql"
	"testing"

	"product-api/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	p := &model.Product{
		Name:        "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Price:       99.99,
	}

	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs(p.Name, p.Category, p.Description, p.Price).
		WillReturnRows(sqlmock.NewRows([]string{"product_id"}).AddRow("1"))

	err = model.CreateProduct(db, p)
	assert.NoError(t, err)
	assert.Equal(t, "1", p.ProductID)
}

func TestCreateProductError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	p := &model.Product{
		Name:        "Test Product",
		Category:    "Test Category",
		Description: "Test Description",
		Price:       99.99,
	}

	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs(p.Name, p.Category, p.Description, p.Price).
		WillReturnError(sql.ErrNoRows)

	err = model.CreateProduct(db, p)
	assert.Error(t, err)
	assert.Equal(t, "", p.ProductID)
}

func TestCreateProductInvalid(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	p := &model.Product{
		Name:        "",
		Category:    "Test Category",
		Description: "Test Description",
		Price:       99.99,
	}

	mock.ExpectQuery(`INSERT INTO products`).
		WithArgs(p.Name, p.Category, p.Description, p.Price).
		WillReturnError(sql.ErrNoRows)

	err = model.CreateProduct(db, p)
	assert.Error(t, err)
	assert.Equal(t, "", p.ProductID)
}

func TestGetTotalProductsCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM products`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))

	total, err := model.GetTotalProductsCount(db)
	assert.NoError(t, err)
	assert.Equal(t, 10, total)
}

func TestFetchProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products`).
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "category", "description", "price"}).
			AddRow("1", "Test Product 1", "Category 1", "Description 1", 10.00).
			AddRow("2", "Test Product 2", "Category 2", "Description 2", 20.00))

	products, err := model.FetchProducts(db, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, "Test Product 1", products[0].Name)
	assert.Equal(t, "Test Product 2", products[1].Name)
}

func TestFetchProductByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 1

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products WHERE product_id = \$1`).
		WithArgs(productID).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "category", "description", "price"}).
			AddRow("1", "Test Product", "Category", "Description", 10.00))

	product, err := model.FetchProductByID(db, productID)
	assert.NoError(t, err)
	assert.Equal(t, "Test Product", product.Name)
}

func TestFetchProductByIDNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 999

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products WHERE product_id = \$1`).
		WithArgs(productID).
		WillReturnError(sql.ErrNoRows)

	product, err := model.FetchProductByID(db, productID)
	assert.Error(t, err)
	assert.Nil(t, product)
}

func TestFetchProductsEmpty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products`).
		WithArgs(limit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"product_id", "name", "category", "description", "price"}))

	products, err := model.FetchProducts(db, limit, offset)
	assert.NoError(t, err)
	assert.Len(t, products, 0)
}

func TestFetchProductsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	limit := 10
	offset := 0

	mock.ExpectQuery(`SELECT product_id, name, category, description, price FROM products`).
		WithArgs(limit, offset).
		WillReturnError(sql.ErrConnDone)

	products, err := model.FetchProducts(db, limit, offset)
	assert.Error(t, err)
	assert.Nil(t, products)
}

func TestProductUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 1
	p := &model.Product{
		Name:        "Updated Product",
		Category:    "Updated Category",
		Description: "Updated Description",
		Price:       99.99,
	}

	mock.ExpectExec(`UPDATE products SET name = \$1, category = \$2, description = \$3, price = \$4 WHERE product_id = \$5`).
		WithArgs(p.Name, p.Category, p.Description, p.Price, productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rowsAffected, err := model.UpdateProduct(db, productID, p)
	assert.NoError(t, err)
	assert.Equal(t, true, rowsAffected)
	assert.NoError(t, err)
}

func TestProductUpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 1
	p := &model.Product{
		Name:        "Updated Product",
		Category:    "Updated Category",
		Description: "Updated Description",
		Price:       99.99,
	}

	mock.ExpectExec(`UPDATE products SET name = \$1, category = \$2, description = \$3, price = \$4 WHERE product_id = \$5`).
		WithArgs(p.Name, p.Category, p.Description, p.Price, productID).
		WillReturnError(sql.ErrNoRows)

	rowsAffected, err := model.UpdateProduct(db, productID, p)
	assert.Error(t, err)
	assert.Equal(t, false, rowsAffected)
}

func TestProductDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 1

	mock.ExpectExec(`DELETE FROM products WHERE product_id = \$1`).
		WithArgs(productID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rowsAffected, err := model.DeleteProduct(db, productID)
	assert.NoError(t, err)
	assert.Equal(t, true, rowsAffected)
}

func TestProductDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 1

	mock.ExpectExec(`DELETE FROM products WHERE product_id = \$1`).
		WithArgs(productID).
		WillReturnError(sql.ErrNoRows)

	rowsAffected, err := model.DeleteProduct(db, productID)
	assert.Error(t, err)
	assert.Equal(t, false, rowsAffected)
	assert.Error(t, err)
}

func TestProductDeleteNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	productID := 1

	mock.ExpectExec(`DELETE FROM products WHERE product_id = \$1`).
		WithArgs(productID).
		WillReturnError(sql.ErrNoRows)

	rowsAffected, err := model.DeleteProduct(db, productID)
	assert.Error(t, err)
	assert.Equal(t, false, rowsAffected)
}
