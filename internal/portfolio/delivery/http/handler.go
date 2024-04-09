package http

import (
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/internal/portfolio"
	"net/http"
)

type Handler struct {
	logger  logging.Logger
	useCase portfolio.UseCase
}

func NewHandler(logger logging.Logger, useCase portfolio.UseCase) *Handler {
	return &Handler{
		logger:  logger,
		useCase: useCase,
	}
}

func (h *Handler) Create(C *gin.Context) {

}

func (h *Handler) Update(c *gin.Context) {
	err := h.useCase.Update(c.Request.Context(), 0)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

/*
func (h *Handler) Test(c *gin.Context) {
	h.logger.Debug("Test")
	err := h.useCase.Ha(c, h.logger)
	h.logger.Debug("", err)
	c.AbortWithStatus(http.StatusAccepted)
}
*/
