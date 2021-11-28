package service

import (
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/repository"
	"github.com/go-playground/validator/v10"
)

type Balance interface {
	Refill(input *model.MoveMoneyModel) (int, error)
	Withdraw(input *model.MoveMoneyModel) (int, error)
	GetByCustomerId(customerId int) (internal.Balance, error)
	create(customerId int, amount int) (int, error)
}

type Service struct {
	Balance
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Balance: NewBalanceService(repo, validator.New()),
	}
}
