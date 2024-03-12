package server

import (
	"go-portfolios-tracker/internal/auth"
	authusecase "go-portfolios-tracker/internal/auth/usecase"
	"net/http"
)

type App struct {
	httpServer *http.Server
	authUC     auth.UseCase
}

func NewApp() *App {
	return &App{
		db :=
		authUC: authusecase.NewAuthUseCase(),
	}
}
