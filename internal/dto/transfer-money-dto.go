package dto

import "github.com/gavrl/app/internal/model"

type TransferMoneyDTO struct {
	CustomerIdFrom int               `json:"customer_from" binding:"required"`
	CustomerIdTo   int               `json:"customer_to" binding:"required"`
	Amount         model.PennyAmount `json:"amount" binding:"required"`
}
