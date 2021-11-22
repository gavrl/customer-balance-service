package handler

import (
	"github.com/gavrl/app/internal/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) refill(c *gin.Context) {
	var input model.RefillDto

	if err := c.BindJSON(&input); err != nil {
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
