package repositories

import (
	"errors"
	"strings"

	"github.com/CAVAh/api-tech-challenge/src/core/domain/entities"
	"github.com/CAVAh/api-tech-challenge/src/infra/db/database"
	"github.com/CAVAh/api-tech-challenge/src/infra/db/models"
)

type CustomerRepository struct {
	DB database.Database
}

func (r CustomerRepository) Create(entity *entities.Customer) (*entities.Customer, error) {
	customer := models.Customer{
		Name:  entity.Name,
		CPF:   entity.CPF,
		Email: entity.Email,
	}

	if err := r.DB.Create(&customer); err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, errors.New("cliente já existe no sistema")
		} else {
			return nil, errors.New("ocorreu um erro desconhecido ao criar o cliente")
		}
	}

	result := customer.ToDomain()

	return &result, nil
}

func (r CustomerRepository) FindFirstByCpf(entity *entities.Customer) (*entities.Customer, error) {
	var customer models.Customer

	db := r.DB.Where("cpf = ?", entity.CPF)
	err := db.First(&customer).Error

	if err != nil {
		return nil, err
	}

	result := customer.ToDomain()

	return &result, nil
}
