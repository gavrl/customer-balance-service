package model

type RefillDto struct {
	CustomerId int     `json:"customer_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
}
