package usecase_test

import (
	"testing"

	"product-service/internal/domain"
	"product-service/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(product *domain.Product) (string, error) {
	args := m.Called(product)
	return args.String(0), args.Error(1)
}

func (m *MockRepo) GetProductsByCategory(category string) ([]*domain.Product, error) {
	args := m.Called(category)
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func (m *MockRepo) GetByID(id string) (*domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockRepo) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepo) List() ([]*domain.Product, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Product), args.Error(1)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := new(MockRepo)
	product := &domain.Product{
		Name:        "Test Product",
		Description: "Test Desc",
		Category:    "Test",
		Price:       99.9,
		Quantity:    10,
		UserID:      "123",
	}
	mockRepo.On("List").Return([]*domain.Product{}, nil) // для инициализации кэша
	mockRepo.On("Create", product).Return("abc123", nil)

	use := usecase.NewProductUsecase(mockRepo)
	id, err := use.Create(product)

	assert.NoError(t, err)
	assert.Equal(t, "abc123", id)
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRepo.On("List").Return([]*domain.Product{}, nil)
	mockRepo.On("Delete", "del123").Return(nil)

	use := usecase.NewProductUsecase(mockRepo)
	err := use.Delete("del123")

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Delete", "del123")
}

func TestGetByID_FromRepo(t *testing.T) {
	mockRepo := new(MockRepo)
	expected := &domain.Product{ID: "id123", Name: "Test"}
	mockRepo.On("List").Return([]*domain.Product{}, nil)
	mockRepo.On("GetByID", "id123").Return(expected, nil)

	use := usecase.NewProductUsecase(mockRepo)
	product, err := use.GetByID("id123")

	assert.NoError(t, err)
	assert.Equal(t, "Test", product.Name)
}

func TestListProducts(t *testing.T) {
	mockRepo := new(MockRepo)
	expected := []*domain.Product{
		{ID: "1", Name: "One"},
		{ID: "2", Name: "Two"},
	}
	mockRepo.On("List").Return(expected, nil)

	use := usecase.NewProductUsecase(mockRepo)
	result, err := use.List()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
