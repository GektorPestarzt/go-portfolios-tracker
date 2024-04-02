package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"go-portfolios-tracker/internal/auth"
	"go-portfolios-tracker/internal/auth/usecase"
	"go-portfolios-tracker/internal/logging"
	"net/http"
)

type Handler struct {
	logger  *logging.Logger
	useCase *usecase.AuthUseCase
}

func NewHandler(logger *logging.Logger, useCase *usecase.AuthUseCase) *Handler {
	return &Handler{
		logger:  logger,
		useCase: useCase,
	}
}

func (h *Handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is list of users"))
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) SignUp(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := h.useCase.SignUp(c.Request.Context(), inp.Username, inp.Password); err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	inp := new(signInput)

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp.Username, inp.Password)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, signInResponse{Token: token})
}

/*
func (h *Handler) GetUserByUUID(c *gin.Context) {

}
*/
/*
func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("this is update auth"))
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("this is partially update auth"))
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(204)
	w.Write([]byte("this is delete auth"))
}
*/
