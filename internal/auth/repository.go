package auth

import (
	"context"
	"go-portfolios-tracker/internal/models"
)

type Repository interface {
	Get(ctx context.Context, username, password string) (*models.User, error)
	Add(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, uuid int) error
}
