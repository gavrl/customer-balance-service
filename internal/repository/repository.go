package repository

import (
	"github.com/gavrl/app/internal"
	"github.com/jmoiron/sqlx"
)

const (
	balanceTable = "balance"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type Balance interface {
	RefillBalance(balance internal.Balance) (int, error)
	GetBalanceByCustomerId(customerId int) (internal.Balance, error)
}

type Repository struct {
	Balance
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Balance: NewBalancePostgres(db),
	}
}
