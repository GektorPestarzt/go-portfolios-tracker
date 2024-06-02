package account

import (
	"context"
	"go-portfolios-tracker/internal/models"
)

type Repository interface {
	Init(ctx context.Context, username string, token string, typeA string) (int64, error)
	GetToken(ctx context.Context, id int64) (string, error)
	Update(ctx context.Context, account *models.Account) error
	UpdateStatus(ctx context.Context, id int64, status models.Status) error
	Get(ctx context.Context, id int64) (*models.Account, error)
	Delete(ctx context.Context, id int64) error
}
