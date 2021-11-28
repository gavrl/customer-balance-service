package handler

import (
	"github.com/gavrl/app/internal/service"
	"github.com/gavrl/app/pkg/formatter"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services  *service.Service
	formatter *formatter.JSONFormatter
}

func NewHandler(services *service.Service, frmt *formatter.JSONFormatter) *Handler {
	return &Handler{
		services:  services,
		formatter: frmt,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	balance := router.Group("/balance")
	{
		balance.POST("/refill", h.refill)
		balance.POST("/withdraw", h.withdraw)
		balance.GET("/:id", h.getBalanceByCustomerId)
	}

	return router
}
