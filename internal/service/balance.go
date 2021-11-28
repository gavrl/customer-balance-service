package service

import (
	"database/sql"
	"errors"
	"github.com/gavrl/app/internal"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/repository"
	"github.com/go-playground/validator/v10"
)

type BalanceService struct {
	repo      repository.BalanceRepository
	validator *validator.Validate
}

func NewBalanceService(repo repository.BalanceRepository, validator *validator.Validate) *BalanceService {
	return &BalanceService{repo: repo, validator: validator}
}

func (bs BalanceService) Refill(moveMoneyModel *model.MoveMoneyModel) (int, error) {
	var amount int

	verr := bs.validateMoveMoneyModel(moveMoneyModel)
	if verr != nil {
		return 0, verr
	}

	balance, err := bs.GetByCustomerId(moveMoneyModel.CustomerId)
	if err != nil {
		// if balance not found, create and return amount
		if errors.Is(err, NotExistsCustomerError{CustomerId: moveMoneyModel.CustomerId}) {
			amount, err = bs.create(moveMoneyModel.CustomerId, moveMoneyModel.Amount)
			if err != nil {
				return 0, err
			}
			return amount, err
		} else {
			return 0, err
		}
	}

	balance.Amount += moveMoneyModel.Amount

	err = bs.updateAmount(balance)
	if err != nil {
		return 0, err
	}

	return balance.Amount, nil
}

func (bs BalanceService) Withdraw(moveMoneyModel *model.MoveMoneyModel) (int, error) {
	verr := bs.validateMoveMoneyModel(moveMoneyModel)
	if verr != nil {
		return 0, verr
	}

	balance, err := bs.GetByCustomerId(moveMoneyModel.CustomerId)
	if err != nil {
		return 0, err
	}

	if balance.Amount-moveMoneyModel.Amount < 0 {
		return 0, NotEnoughFundsError{}
	}

	balance.Amount -= moveMoneyModel.Amount

	err = bs.updateAmount(balance)
	if err != nil {
		return 0, err
	}

	return balance.Amount, nil
}

func (bs BalanceService) GetByCustomerId(customerId int) (internal.Balance, error) {
	balance, err := bs.repo.GetByCustomerId(customerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return balance, NotExistsCustomerError{CustomerId: customerId}
		} else {
			return balance, err
		}
	}
	return balance, nil
}

func (bs BalanceService) Transfer(transferModel *model.TransferMoneyModel) error {
	err := bs.validator.Struct(transferModel)
	if err != nil {
		return err.(validator.ValidationErrors)
	}

	balanceFrom, err := bs.GetByCustomerId(transferModel.CustomerIdFrom)
	if err != nil {
		return err
	}
	balanceTo, err := bs.GetByCustomerId(transferModel.CustomerIdTo)
	if err != nil {
		return err
	}

	if balanceFrom.Amount-transferModel.Amount < 0 {
		return NotEnoughFundsError{}
	}

	balanceFrom.Amount -= transferModel.Amount
	balanceTo.Amount += transferModel.Amount

	err = bs.repo.TransferMoney(&balanceFrom, &balanceTo)
	if err != nil {
		return err
	}

	return nil
}

func (bs BalanceService) create(customerId int, amount int) (int, error) {
	amnt, err := bs.repo.Create(customerId, amount)
	if err != nil {
		return 0, err
	}
	return amnt, nil
}

func (bs BalanceService) updateAmount(balance internal.Balance) error {
	err := bs.repo.UpdateAmount(&balance)
	if err != nil {
		return err
	}
	return nil
}

func (bs BalanceService) validateMoveMoneyModel(model *model.MoveMoneyModel) error {
	err := bs.validator.Struct(model)
	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
