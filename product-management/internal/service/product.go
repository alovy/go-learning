package service

import (
	"database/sql"
	"product-api/internal/model"
)

type productService struct {
	db *sql.DB
}

func NewProductService(db *sql.DB) ProductService {
	return &productService{db: db}
}

func (s *productService) CreateProduct(p *model.Product) error {
	return model.CreateProduct(s.db, p)
}

func (s *productService) GetProductByID(id int) (*model.Product, error) {
	return model.FetchProductByID(s.db, id)
}

func (s *productService) GetAllProducts(limit, offset int) ([]model.Product, int, error) {
	total, err := model.GetTotalProductsCount(s.db)
	if err != nil {
		return nil, 0, err
	}

	products, err := model.FetchProducts(s.db, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (s *productService) UpdateProduct(id int, p *model.Product) (bool, error) {
	return model.UpdateProduct(s.db, id, p)
}

func (s *productService) DeleteProduct(id int) (bool, error) {
	return model.DeleteProduct(s.db, id)
}
