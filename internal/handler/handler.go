package handler

import (
	"github.com/gavrl/app/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	balance := router.Group("/balance")
	{
		balance.POST("/refill", h.refill)
	}

	return router
}
