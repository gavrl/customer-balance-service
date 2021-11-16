package repository

import (
	"github.com/gavrl/app/internal"
	"github.com/jmoiron/sqlx"
)

type BalancePostgres struct {
	db *sqlx.DB
}

func (a BalancePostgres) RefillBalance(balance internal.Balance) (int, error) {
	panic("implement me")
}

func (a BalancePostgres) GetBalanceByCustomerId(customerId int) (internal.Balance, error) {
	panic("implement me")
}

func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}
