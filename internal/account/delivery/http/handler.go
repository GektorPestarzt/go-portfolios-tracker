package http

import (
	"encoding/json"
	"go-portfolios-tracker/internal/account"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/logging"
	"go-portfolios-tracker/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	id, err := h.useCase.Init(c.Request.Context(), c.Value(auth.CtxUserKey).(string), inp.Token)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{"id": strconv.Itoa(int(id))})
	c.Status(http.StatusCreated)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.useCase.UpdateStatus(c.Request.Context(), int64(id), models.Process)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// h.broker.Publish(c.Request.Context(), fmt.Sprintf("%d", id))
	// c.Status(http.StatusCreated)

	err = h.useCase.Update(c.Request.Context(), int64(id))
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

	acc, err := h.useCase.Get(c.Request.Context(), int64(id))
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
