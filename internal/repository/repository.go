package repository

import (
	"github.com/gavrl/app/internal"
	"github.com/jmoiron/sqlx"
)

const balanceTable = "balance"

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type BalanceRepository interface {
	GetByCustomerId(customerId int) (internal.Balance, error)
	Create(customerId int, amount int) (int, error)
	UpdateAmount(balance *internal.Balance) error
	TransferMoney(balanceFrom *internal.Balance, balanceTo *internal.Balance) error
}

type Repository struct {
	BalanceRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		BalanceRepository: NewBalancePostgres(db),
	}
}
