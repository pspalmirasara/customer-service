package usecases

import (
	"fmt"
	"testing"

	"github.com/CAVAh/api-tech-challenge/src/adapters/gateways"
	"github.com/CAVAh/api-tech-challenge/src/core/domain/dtos"
	"github.com/CAVAh/api-tech-challenge/src/core/domain/entities"
	"github.com/stretchr/testify/assert"
)

type mockListCustomerRepository struct {
	gateways.CustomerRepository
	mockFindFirstByCpf func(*entities.Customer) (*entities.Customer, error)
}

func (m *mockListCustomerRepository) FindFirstByCpf(customer *entities.Customer) (*entities.Customer, error) {
	return m.mockFindFirstByCpf(customer)
}

func TestListCustomerUsecase_Execute(t *testing.T) {
	mockCustomerRepo := &mockListCustomerRepository{}

	usecase := ListCustomerUsecase{
		CustomerRepository: mockCustomerRepo,
	}

	inputDto := dtos.ListCustomerDto{
		CPF: "12345678900",
	}

	t.Run("valid input", func(t *testing.T) {
		mockCustomerRepo.mockFindFirstByCpf = func(customer *entities.Customer) (*entities.Customer, error) {
			return &entities.Customer{}, nil
		}

		_, err := usecase.Execute(inputDto)
		assert.NoError(t, err)
	})

	t.Run("erro ao listar cliente", func(t *testing.T) {
		mockCustomerRepo.mockFindFirstByCpf = func(customer *entities.Customer) (*entities.Customer, error) {
			return nil, fmt.Errorf("Erro ao listar cliente")
		}

		_, err := usecase.Execute(inputDto)
		assert.NoError(t, err)
	})

	t.Run("cpf vazio", func(t *testing.T) {
		inputDto.CPF = ""
		_, err := usecase.Execute(inputDto)
		assert.NoError(t, err)
	})
}
