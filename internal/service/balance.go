package service

import (
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/repository"
)

type BalanceService struct {
	repo repository.Balance
}

func NewBalanceService(repo repository.Balance) *BalanceService {
	return &BalanceService{repo: repo}
}

func (b BalanceService) RefillBalance(input model.RefillDto) (int, error) {
	panic("implement me")
}

func (b BalanceService) GetBalanceByCustomerId(customerId int) (internal.Balance, error) {
	panic("implement me")
}
