package handler

import (
	"errors"
	"github.com/gavrl/app/internal/dto"
	"github.com/gavrl/app/internal/model"
	"github.com/gavrl/app/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
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

func (h *Handler) getBalanceByCustomerId(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	balance, err := h.services.Balance.GetByCustomerId(id)
	if err != nil {
		if errors.Is(err, service.NotExistsCustomerError{CustomerId: id}) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
		} else {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
		}
		return
	}

	outDTO := dto.MoveMoneyDTO{
		CustomerId: balance.CustomerId,
		Amount:     model.PennyAmount{Int: balance.Amount},
	}

	c.JSON(http.StatusOK, &outDTO)
}

func (h *Handler) transfer(c *gin.Context) {
	var input dto.TransferMoneyDTO
	var verr validator.ValidationErrors

	if err := c.ShouldBind(&input); err != nil {
		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Balance.Transfer(&model.TransferMoneyModel{
		CustomerIdFrom: input.CustomerIdFrom,
		CustomerIdTo:   input.CustomerIdTo,
		Amount:         input.Amount.Int,
	})

	if err != nil {
		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
		} else {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())

		}
		return
	}
}
