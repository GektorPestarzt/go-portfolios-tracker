package http

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go-portfolios-tracker/internal/account"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/logging"
	"net/http"
	"strconv"
)

type Handler struct {
	logger  logging.Logger
	useCase account.UseCase
}

func NewHandler(logger logging.Logger, useCase account.UseCase) *Handler {
	return &Handler{
		logger:  logger,
		useCase: useCase,
	}
}

type createInput struct {
	Token string `json:"token"`
}

func (h *Handler) Create(c *gin.Context) {
	inp := new(createInput)
	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.useCase.Init(c.Request.Context(), c.Value(auth.CtxUserKey).(string), inp.Token); err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.useCase.Update(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	acc, err := h.useCase.Get(c.Request.Context(), id)
	b, err := json.Marshal(acc)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, string(b))
}

/*
func (h *Handler) Test(c *gin.Context) {
	h.logger.Debug("Test")
	err := h.useCase.Ha(c, h.logger)
	h.logger.Debug("", err)
	c.AbortWithStatus(http.StatusAccepted)
}
*/
