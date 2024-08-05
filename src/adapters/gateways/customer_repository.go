package gateways

import (
	"github.com/CAVAh/api-tech-challenge/src/core/domain/entities"
)

type CustomerRepository interface {
	Create(customer *entities.Customer) (*entities.Customer, error)
	FindFirstByCpf(customer *entities.Customer) (*entities.Customer, error)
}
