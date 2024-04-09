package http

import (
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/internal/portfolio"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, useCase portfolio.UseCase, logger logging.Logger) {
	h := NewHandler(logger, useCase)

	portfolios := router.Group("/portfolios")
	{
		// portfolios.GET("/test", h.Test)
		portfolios.POST("", h.Create)
		portfolios.PUT("", h.Update)
	}
}
