package portfolio

import (
	"context"
	"go-portfolios-tracker/internal/models"
)

type UseCase interface {
	Init(ctx context.Context, token string) (int, error)
	Update(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Account, error)
	Delete(ctx context.Context, id int) error
	// Ha(ctx context.Context, logger logging.Logger) error
}
