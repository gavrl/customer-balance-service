package service

import (
	"database/sql"
	"errors"
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/repository"
)

type BalanceService struct {
	repo repository.BalanceRepository
}

func NewBalanceService(repo repository.BalanceRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (bs BalanceService) Refill(input model.RefillDto) (int, error) {
	var amount int
	balance, err := bs.GetByCustomerId(input.CustomerId)
	if err != nil {
		// if balance not found, create and return amount
		if errors.Is(err, sql.ErrNoRows) {
			amount, err = bs.create(input.CustomerId, input.Amount.Int)
			if err != nil {
				return 0, err
			}
			return amount, err
		} else {
			return 0, err
		}
	}

	balance.Amount += input.Amount.Int

	err = bs.updateAmount(balance)
	if err != nil {
		return 0, err
	}

	return balance.Amount, nil
}

func (bs BalanceService) GetByCustomerId(customerId int) (internal.Balance, error) {
	var balance internal.Balance
	balance, err := bs.repo.GetByCustomerId(customerId)
	if err != nil {
		return balance, err
	}
	return balance, nil
}

func (bs BalanceService) create(customerId int, amount int) (int, error) {
	amnt, err := bs.repo.Create(customerId, amount)
	if err != nil {
		return 0, err
	}
	return amnt, nil
}

func (bs BalanceService) updateAmount(balance internal.Balance) error {
	err := bs.repo.UpdateAmount(balance)
	if err != nil {
		return err
	}
	return nil
}
