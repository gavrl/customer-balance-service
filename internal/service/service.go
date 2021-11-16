package service

import (
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/repository"
)

type Balance interface {
	RefillBalance(input model.RefillDto) (int, error)
	GetBalanceByCustomerId(customerId int) (internal.Balance, error)
}

type Service struct {
	Balance
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Balance: NewBalanceService(repo),
	}
}
