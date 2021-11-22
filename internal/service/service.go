package service

import (
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/repository"
)

type Balance interface {
	Refill(input model.RefillDto) (int, error)
	GetByCustomerId(customerId int) (internal.Balance, error)
	create(customerId int, amount int) (int, error)
}

type Service struct {
	Balance
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Balance: NewBalanceService(repo),
	}
}
