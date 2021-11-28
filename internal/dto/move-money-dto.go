package dto

import "github.com/gavrl/app/internal/model"

type MoveMoneyDTO struct {
	CustomerId int               `json:"customer_id" binding:"required"`
	Amount     model.PennyAmount `json:"amount" binding:"required"`
}
