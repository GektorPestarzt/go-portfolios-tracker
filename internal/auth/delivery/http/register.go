package http

import (
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/logging"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "go-portfolios-tracker/docs"
)

func RegisterHTTPEndpoints(router *gin.Engine, useCase auth.UseCase, logger logging.Logger) {
	h := NewHandler(useCase, logger)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authEndpoints := router.Group("/auth")
	{
		authEndpoints.POST("/sign-up", h.SignUp)
		authEndpoints.GET("/sign-in", h.SignIn)
	}
}
