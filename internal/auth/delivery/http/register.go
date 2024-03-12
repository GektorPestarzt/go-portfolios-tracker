package http

import (
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/auth/usecase"
	"go-portfolios-tracker/internal/logging"
)

func RegisterHTTPEndpoints(router *gin.Engine, logger *logging.Logger, useCase *usecase.AuthUseCase) {
	h := NewHandler(logger, useCase)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.GET("/sign-in", h.SignIn)
	}
}
