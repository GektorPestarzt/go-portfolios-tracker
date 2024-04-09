package http

import (
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/logging"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	usecase auth.UseCase
	logger  logging.Logger
}

func NewAuthMiddleware(usecase auth.UseCase, logger logging.Logger) gin.HandlerFunc {
	return (&AuthMiddleware{
		usecase: usecase,
		logger:  logger,
	}).Handle
}

func (m *AuthMiddleware) Handle(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if headerParts[0] != "Bearer" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	username, err := m.usecase.ParseToken(c.Request.Context(), headerParts[1])
	if err != nil {
		if err == auth.ErrInvalidAccessToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		m.logger.Debug("%v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Set(auth.CtxUserKey, username)
}
