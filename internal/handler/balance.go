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

	_, err := h.services.Balance.RefillBalance(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "XE-XE-XE",
	})
}
