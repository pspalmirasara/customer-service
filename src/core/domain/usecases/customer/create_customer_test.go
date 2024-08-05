package usecases

import (
	"fmt"
	"testing"

	"github.com/CAVAh/api-tech-challenge/src/adapters/gateways"
	"github.com/CAVAh/api-tech-challenge/src/core/domain/dtos"
	"github.com/CAVAh/api-tech-challenge/src/core/domain/entities"
	"github.com/stretchr/testify/assert"
)

type mockCreateCustomerRepository struct {
	gateways.CustomerRepository
	mockCreate func(*entities.Customer) (*entities.Customer, error)
}

func (m *mockCreateCustomerRepository) Create(customer *entities.Customer) (*entities.Customer, error) {
	return m.mockCreate(customer)
}

func TestCreateCustomerUsecase_Execute(t *testing.T) {
	mockCustomerRepo := &mockCreateCustomerRepository{}

	usecase := CreateCustomerUsecase{
		CustomerRepository: mockCustomerRepo,
	}

	inputDto := dtos.CreateCustomerDto{
		Name:  "John Doe",
		CPF:   "12345678900",
		Email: "john@example.com",
	}

	t.Run("valid input", func(t *testing.T) {
		mockCustomerRepo.mockCreate = func(customer *entities.Customer) (*entities.Customer, error) {
			return &entities.Customer{}, nil
		}

		_, err := usecase.Execute(inputDto)
		assert.NoError(t, err)
	})

	t.Run("erro ao criar cliente", func(t *testing.T) {
		mockCustomerRepo.mockCreate = func(customer *entities.Customer) (*entities.Customer, error) {
			return nil, fmt.Errorf("Erro ao criar cliente")
		}

		_, err := usecase.Execute(inputDto)
		assert.Error(t, err)
	})
}
