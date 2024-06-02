package http

import (
	"go-portfolios-tracker/internal/account"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/pkg/broker"

	"github.com/gin-gonic/gin"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, useCase account.UseCase, broker broker.Broker, logger logging.Logger) {
	h := NewHandler(logger, useCase, broker)

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
