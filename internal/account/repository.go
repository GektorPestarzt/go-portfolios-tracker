package account

import (
	"context"
	"go-portfolios-tracker/internal/models"
)

type Repository interface {
	Init(ctx context.Context, username string, token string, typeA string) error
	GetToken(ctx context.Context, id int) (string, error)
	Update(ctx context.Context, account *models.Account) error
	Get(ctx context.Context, id int) (*models.Account, error)
	Delete(ctx context.Context, id int) error
}
