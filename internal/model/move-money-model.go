package model

type MoveMoneyModel struct {
	CustomerId int `json:"customer_id" validate:"gt=0"`
	Amount     int `json:"amount" validate:"gt=0"`
}
