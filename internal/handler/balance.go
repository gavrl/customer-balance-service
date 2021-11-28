package handler

import (
	"errors"
	"github.com/gavrl/app/internal/dto"
	"github.com/gavrl/app/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) refill(c *gin.Context) {
	var input dto.MoveMoneyDTO
	var verr validator.ValidationErrors

	if err := c.ShouldBind(&input); err != nil {

		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := h.services.Balance.Refill(&model.MoveMoneyModel{
		CustomerId: input.CustomerId,
		Amount:     input.Amount.Int,
	})

	if err != nil {
		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())

		}
		return
	}

	penny := model.PennyAmount{Int: amount}

	c.JSON(http.StatusOK, map[string]interface{}{
		"amount": &penny,
	})
}

func (h *Handler) withdraw(c *gin.Context) {
	var input dto.MoveMoneyDTO
	var verr validator.ValidationErrors

	if err := c.ShouldBind(&input); err != nil {
		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := h.services.Balance.Withdraw(&model.MoveMoneyModel{
		CustomerId: input.CustomerId,
		Amount:     input.Amount.Int,
	})

	if err != nil {
		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())

		}
		return
	}

	penny := model.PennyAmount{Int: amount}

	c.JSON(http.StatusOK, map[string]interface{}{
		"amount": &penny,
	})
}
