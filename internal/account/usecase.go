package account

import (
	"context"
	"go-portfolios-tracker/internal/models"
)

type UseCase interface {
	Init(ctx context.Context, username string, token string) (int64, error)
	Update(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status models.Status) error
	Get(ctx context.Context, id int64) (*models.Account, error)
	Delete(ctx context.Context, id int64) error
	// Ha(ctx context.Context, logger logging.Logger) error
}
