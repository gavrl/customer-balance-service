package repository

import (
	"fmt"
	"github.com/gavrl/app/internal"
	"github.com/jmoiron/sqlx"
)

type BalancePostgres struct {
	db *sqlx.DB
}

func NewBalancePostgres(db *sqlx.DB) *BalancePostgres {
	return &BalancePostgres{db: db}
}

func (r BalancePostgres) Refill(balance internal.Balance) error {
	panic("implement me")
}

func (r BalancePostgres) GetByCustomerId(customerId int) (internal.Balance, error) {
	var balance internal.Balance

	query := fmt.Sprintf("SELECT id, customer_id, amount FROM %s where customer_id = $1",
		balanceTable)
	err := r.db.Get(&balance, query, customerId)
	if err != nil {
		return balance, err
	}

	return balance, nil
}

func (r BalancePostgres) Create(customerId int, amount int) (int, error) {
	var amnt int
	createListQuery := fmt.Sprintf("INSERT INTO %s (customer_id, amount) VALUES ($1, $2) RETURNING amount", balanceTable)
	row := r.db.QueryRow(createListQuery, customerId, amount)
	if err := row.Scan(&amnt); err != nil {
		return 0, err
	}

	return amnt, nil
}

func (r BalancePostgres) UpdateAmount(balance internal.Balance) error {
	query := fmt.Sprintf("UPDATE %s SET amount = %d WHERE id = %d",
		balanceTable, balance.Amount, balance.Id)

	_, err := r.db.Exec(query)
	return err
}
