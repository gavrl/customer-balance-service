package handler

import (
	"errors"
	"github.com/gavrl/app/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func (h *Handler) refill(c *gin.Context) {
	var input model.RefillDto

	if err := c.ShouldBind(&input); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			newValidateErrorResponse(c, h.formatter, verr)
			return
		}

		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	amount, err := h.services.Balance.Refill(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	penny := model.PennyAmount{Int: amount}

	c.JSON(http.StatusOK, map[string]interface{}{
		"amount": &penny,
	})
}
