package http

import (
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/logging"
)

func RegisterHTTPEndpoints(router *gin.Engine, useCase auth.UseCase, logger logging.Logger) {
	h := NewHandler(useCase, logger)

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.GET("/sign-in", h.SignIn)
	}
}
