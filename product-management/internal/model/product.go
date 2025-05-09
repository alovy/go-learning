package model

import (
	"database/sql"
)

type Product struct {
	ProductID   string  `json:"product_id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

func CreateProduct(db *sql.DB, p *Product) error {
	query := `
		INSERT INTO products (name, category, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING product_id`
	return db.QueryRow(query, p.Name, p.Category, p.Description, p.Price).
		Scan(&p.ProductID)
}

func GetTotalProductsCount(db *sql.DB) (int, error) {
	var total int
	err := db.QueryRow(`SELECT COUNT(*) FROM products`).Scan(&total)
	return total, err
}

func FetchProducts(db *sql.DB, limit, offset int) ([]Product, error) {
	query := `
		SELECT product_id, name, category, description, price
		FROM products
		ORDER BY product_id
		LIMIT $1 OFFSET $2`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ProductID, &p.Name, &p.Category, &p.Description, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// FetchProductByID retrieves a product by its ID from the database.
func FetchProductByID(db *sql.DB, productID int) (*Product, error) {
	query := `
		SELECT product_id, name, category, description, price
		FROM products
		WHERE product_id = $1`
	var p Product
	err := db.QueryRow(query, productID).Scan(&p.ProductID, &p.Name, &p.Category, &p.Description, &p.Price)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// UpdateProduct updates an existing product in the database.
func UpdateProduct(db *sql.DB, productID int, p *Product) (bool, error) {
	query := `
		UPDATE products
		SET name = $1, category = $2, description = $3, price = $4
		WHERE product_id = $5`

	res, err := db.Exec(query, p.Name, p.Category, p.Description, p.Price, productID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func DeleteProduct(db *sql.DB, productID int) (bool, error) {
	query := `DELETE FROM products WHERE product_id = $1`
	res, err := db.Exec(query, productID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
