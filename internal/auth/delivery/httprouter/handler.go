package httprouter

import (
	"github.com/julienschmidt/httprouter"
	"go-portfolios-tracker/internal/auth/usecase"
	"go-portfolios-tracker/internal/handlers"
	"go-portfolios-tracker/internal/logging/slog"
	"net/http"
)

var _ handlers.Handler = &handler{}

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	logger  slog.Logger
	useCase usecase.AuthUseCase
}

func NewHandler(logger slog.Logger, useCase usecase.AuthUseCase) handlers.Handler {
	return &handler{
		logger:  logger,
		useCase: useCase,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(usersURL, h.GetList)
	router.POST(usersURL, h.CreateUser)
	router.GET(userURL, h.GetUserByUUID)
	router.PUT(userURL, h.UpdateUser)
	router.PATCH(userURL, h.PartiallyUpdateUser)
	router.DELETE(userURL, h.DeleteUser)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is list of users"))
}

type signInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(201)

	inp := new(signInput)
	//if err := cleanenv.ParseJSON(r, inp); err != nil {
	//}

	if err := h.useCase.SignUp(r.Context(), inp.Username, inp.Password); err != nil {

	}
	w.Write([]byte("this is create auth"))
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.WriteHeader(200)
	w.Write([]byte("this is auth by id"))
}

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
