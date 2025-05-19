package tests

import (
	"product-api/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(p *model.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockProductService) GetProductByID(id int) (*model.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockProductService) GetAllProducts(limit, offset int) ([]model.Product, int, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]model.Product), args.Int(1), args.Error(2)
}

func (m *MockProductService) UpdateProduct(id int, p *model.Product) (bool, error) {
	args := m.Called(id, p)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductService) DeleteProduct(id int) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}
