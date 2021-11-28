package service

import "fmt"

type NotEnoughFunds struct{}

func (e NotEnoughFunds) Error() string {
	return "there are not enough funds"
}

type NotExistsCustomer struct {
	customerId int
}

func (e NotExistsCustomer) Error() string {
	return fmt.Sprintf("customer with id %d does not exists", e.customerId)
}
