package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CAVAh/api-tech-challenge/src/adapters/gateways"
	"github.com/CAVAh/api-tech-challenge/src/core/domain/dtos"
	"github.com/CAVAh/api-tech-challenge/src/core/domain/entities"
	usecases "github.com/CAVAh/api-tech-challenge/src/core/domain/usecases/customer"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do repositório de clientes
type MockCustomerRepository struct {
	gateways.CustomerRepository
	mock.Mock
}

func (m *MockCustomerRepository) Create(customer *entities.Customer) (*entities.Customer, error) {
	args := m.Called(customer)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCustomerRepository) FindFirstByCpf(customer *entities.Customer) (*entities.Customer, error) {
	args := m.Called(customer)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.Customer), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestListCustomer_InvalidInput(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.ListCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.GET("/customers", func(c *gin.Context) {
		ListCustomers(c, &usecase)
	})

	// Criar uma requisição HTTP simulada com parâmetros de query inválidos
	req, _ := http.NewRequest(http.MethodGet, "/customers?cpf=invalid_cpf_format", nil)
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestListCustomer_UseCaseError(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)

	// Configurar o mock para retornar um erro ao listar os clientes
	mockRepo.On("FindFirstByCpf", mock.Anything).Return(nil, errors.New("some error"))

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.ListCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.GET("/customers", func(c *gin.Context) {
		ListCustomers(c, &usecase)
	})

	// Criar uma requisição HTTP simulada com parâmetros de query válidos
	req, _ := http.NewRequest(http.MethodGet, "/customers?cpf=12345678900", nil)
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListCustomers(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)
	expectedCustomer := entities.Customer{
		ID:        1,
		Name:      "Customer 1",
		CPF:       "12345678900",
		Email:     "email@email.com",
		CreatedAt: "2021-01-01",
	}

	// Configurar o mock para retornar o cliente esperado
	mockRepo.On("FindFirstByCpf", mock.Anything).Return(&expectedCustomer, nil)

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.ListCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.GET("/customers", func(c *gin.Context) {
		ListCustomers(c, &usecase)
	})

	// Criar uma requisição HTTP simulada
	req, _ := http.NewRequest(http.MethodGet, "/customers?cpf=12345678900", nil)
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar se o mock foi chamado corretamente
	mockRepo.AssertExpectations(t)
}

func TestListCustomers_WithoutCpf(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.ListCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.GET("/customers", func(c *gin.Context) {
		ListCustomers(c, &usecase)
	})

	// Criar uma requisição HTTP simulada
	req, _ := http.NewRequest(http.MethodGet, "/customers", nil)
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusOK, w.Code)

	// Verificar se o mock não foi chamado
	mockRepo.AssertNotCalled(t, "FindFirstByCpf", mock.Anything)
}

func TestCreateCustomer_InvalidJSON(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.CreateCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.POST("/customers", func(c *gin.Context) {
		CreateCustomer(c, &usecase)
	})

	// Criar uma requisição HTTP simulada com um JSON malformado
	invalidJSON := []byte(`{"name":"New Customer", "cpf":12345678900, "email":"newcustomer@email.com"}`) // CPF como número, deve ser string
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "json: cannot unmarshal")
}

func TestCreateCustomer_InvalidInput(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.CreateCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.POST("/customers", func(c *gin.Context) {
		CreateCustomer(c, &usecase)
	})

	// Criar uma requisição HTTP simulada com um JSON inválido
	invalidInputDto := dtos.CreateCustomerDto{
		Name:  "", // Nome vazio para falhar na validação
		CPF:   "12345678900",
		Email: "invalid-email", // Email inválido
	}
	invalidInputJSON, _ := json.Marshal(invalidInputDto)
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(invalidInputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestCreateCustomer_UseCaseError(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)

	// Configurar o mock para retornar um erro ao criar o cliente
	mockRepo.On("Create", mock.Anything).Return(nil, errors.New("some error"))

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.CreateCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.POST("/customers", func(c *gin.Context) {
		CreateCustomer(c, &usecase)
	})

	// Criar uma requisição HTTP simulada com um corpo JSON
	inputDto := dtos.CreateCustomerDto{
		Name:  "New Customer",
		CPF:   "12345678900",
		Email: "newcustomer@email.com",
	}
	inputJSON, _ := json.Marshal(inputDto)
	req, _ := http.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(inputJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"some error"`)
}

func TestCreateCustomer(t *testing.T) {
	// Configurar o gin em modo de teste
	gin.SetMode(gin.TestMode)

	// Criar o mock do repositório de clientes
	mockRepo := new(MockCustomerRepository)
	expectedCustomer := entities.Customer{
		ID:        1,
		Name:      "Customer 1",
		CPF:       "12345678900",
		Email:     "email@email.com",
		CreatedAt: "2021-01-01",
	}

	// Configurar o mock para retornar o cliente esperado
	mockRepo.On("Create", mock.Anything).Return(&expectedCustomer, nil)

	// Substituir o repositório real pelo mock no usecase
	usecase := usecases.CreateCustomerUsecase{
		CustomerRepository: mockRepo,
	}

	// Configurar o controlador com o mock do usecase
	r := gin.Default()
	r.POST("/customers", func(c *gin.Context) {
		CreateCustomer(c, &usecase)
	})

	// Criar uma string JSON com os dados de entrada
	inputJSON := `{"name":"Customer 1","cpf":"12345678900","email":"email@email.com"}`

	// Converter a string JSON em um buffer
	reqBody := bytes.NewBufferString(inputJSON)

	// Criar uma requisição HTTP simulada
	req, _ := http.NewRequest(http.MethodPost, "/customers", reqBody)
	req.Header.Set("Content-Type", "application/json") // Definir o cabeçalho da requisição para indicar o tipo de conteúdo

	w := httptest.NewRecorder()

	// Executar a requisição
	r.ServeHTTP(w, req)

	// Verificar o resultado
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.JSONEq(t, `{"cpf":"12345678900", "createdAt":"2021-01-01", "email":"email@email.com", "id":1, "name":"Customer 1"}`, w.Body.String())

	// Verificar se o mock foi chamado corretamente
	mockRepo.AssertExpectations(t)
}
