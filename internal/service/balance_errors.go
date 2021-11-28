package service

import "fmt"

type NotEnoughFundsError struct{}

func (e NotEnoughFundsError) Error() string {
	return "there are not enough funds"
}

type NotExistsCustomerError struct {
	CustomerId int
}

func (e NotExistsCustomerError) Error() string {
	return fmt.Sprintf("customer with id %d does not exists", e.CustomerId)
}
