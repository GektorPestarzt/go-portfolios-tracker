package http

import (
	"go-portfolios-tracker/internal/account"
	"go-portfolios-tracker/internal/logging"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, useCase account.UseCase, logger logging.Logger) {
	h := NewHandler(logger, useCase)

	accounts := router.Group("/accounts")
	{
		tinkoff := accounts.Group("/tinkoff")
		{
			tinkoff.POST("", h.Create)
			tinkoff.GET("/:id", h.Get)
			tinkoff.PUT("/:id", h.Update)
		}
	}
}
