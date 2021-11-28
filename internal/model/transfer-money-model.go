package model

type TransferMoneyModel struct {
	CustomerIdFrom int `json:"customer_from" validate:"gt=0"`
	CustomerIdTo   int `json:"customer_to" validate:"gt=0"`
	Amount         int `json:"amount" validate:"gt=0"`
}
