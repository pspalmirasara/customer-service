package repositories

import (
	"errors"
	"testing"

	"github.com/CAVAh/api-tech-challenge/src/core/domain/entities"
	"github.com/CAVAh/api-tech-challenge/src/infra/db/mocks"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateCustomer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	repo := CustomerRepository{DB: mockDB}

	entity := &entities.Customer{
		Name:  "John Doe",
		CPF:   "12345678901",
		Email: "john@example.com",
	}

	mockDB.EXPECT().Create(gomock.Any()).Return(nil)

	result, err := repo.Create(entity)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "John Doe", result.Name)
}

func TestCreateCustomer_Duplicate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	repo := CustomerRepository{DB: mockDB}

	entity := &entities.Customer{
		Name:  "John Doe",
		CPF:   "12345678901",
		Email: "john@example.com",
	}

	mockDB.EXPECT().Create(gomock.Any()).Return(errors.New("duplicate key value violates unique constraint"))

	result, err := repo.Create(entity)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "cliente já existe no sistema", err.Error())
}

func TestCreateCustomer_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	repo := CustomerRepository{DB: mockDB}

	entity := &entities.Customer{
		Name:  "John Doe",
		CPF:   "12345678901",
		Email: "john@example.com",
	}

	mockDB.EXPECT().Create(gomock.Any()).Return(errors.New("some error"))

	result, err := repo.Create(entity)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "ocorreu um erro desconhecido ao criar o cliente", err.Error())
}

// FIXME: This test is not working
// func TestFindFirstByCpf_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockDB := mocks.NewMockDatabase(ctrl)
// 	repo := CustomerRepository{DB: mockDB}

// 	entity := &entities.Customer{CPF: "12345678901"}
// 	customerModel := models.Customer{
// 		Name:  "John Doe",
// 		CPF:   "12345678901",
// 		Email: "john@example.com",
// 	}

// 	// Mock para retornar a instância correta de gorm.DB
// 	mockDB.EXPECT().Where("cpf = ?", entity.CPF).Return(&gorm.DB{}).Times(1)
// 	mockDB.EXPECT().First(gomock.Any()).DoAndReturn(func(dest interface{}) *gorm.DB {
// 		*dest.(*models.Customer) = customerModel
// 		return &gorm.DB{}
// 	}).Times(1)

// 	result, err := repo.FindFirstByCpf(entity)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, result)
// 	assert.Equal(t, "John Doe", result.Name)
// 	assert.Equal(t, "12345678901", result.CPF)
// 	assert.Equal(t, "john@example.com", result.Email)
// }

// func TestFindFirstByCpf_NotFound(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockDB := mocks.NewMockDatabase(ctrl)
// 	repo := CustomerRepository{DB: mockDB}

// 	entity := &entities.Customer{CPF: "12345678901"}

// 	// Mock the Where and Find calls
// 	mockDB.EXPECT().Where("cpf = ?", entity.CPF).Return(mockDB)
// 	mockDB.EXPECT().Find(gomock.Any()).DoAndReturn(func(dest interface{}) *gorm.DB {
// 		*dest.(*[]models.Customer) = []models.Customer{}
// 		return &gorm.DB{}
// 	})

// 	result, err := repo.FindFirstByCpf(entity)
// 	assert.NoError(t, err)
// 	assert.Nil(t, result)
// }

// func TestFindFirstByCpf_Error(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockDB := mocks.NewMockDatabase(ctrl)
// 	repo := CustomerRepository{DB: mockDB}

// 	entity := &entities.Customer{CPF: "12345678901"}

// 	// Mock the Where and Find calls
// 	mockDB.EXPECT().Where("cpf = ?", entity.CPF).Return(mockDB)
// 	mockDB.EXPECT().Find(gomock.Any()).Return(&gorm.DB{Error: errors.New("database error")})

// 	result, err := repo.FindFirstByCpf(entity)
// 	assert.Error(t, err)
// 	assert.Nil(t, result)
// 	assert.Equal(t, "database error", err.Error())
// }
