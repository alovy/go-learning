package service

import "product-api/internal/model"

type ProductService interface {
	CreateProduct(p *model.Product) error
	GetProductByID(id int) (*model.Product, error)
	GetAllProducts(limit, offset int) ([]model.Product, int, error)
	UpdateProduct(id int, p *model.Product) (bool, error)
	DeleteProduct(id int) (bool, error)
}
