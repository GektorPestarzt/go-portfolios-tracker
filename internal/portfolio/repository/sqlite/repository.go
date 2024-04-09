package sqlite

import (
	"context"
	"database/sql"
	"go-portfolios-tracker/internal/models"
)

type PortfolioRepository struct {
	db *sql.DB
}

func NewPortfolioRepository(db *sql.DB) *PortfolioRepository {
	return &PortfolioRepository{
		db: db,
	}
}

// TODO: remake. Test version
func (pr *PortfolioRepository) Init(ctx context.Context, token string) (int, error) {
	return 0, nil
}

func (pr *PortfolioRepository) GetToken(ctx context.Context, id int) (string, error) {
	return "", nil
}

func (pr *PortfolioRepository) Update(ctx context.Context, account *models.Account) error {
	return nil
}

func (pr *PortfolioRepository) Get(ctx context.Context, id int) (*models.Account, error) {
	return nil, nil
}

func (pr *PortfolioRepository) Delete(ctx context.Context, id int) error {
	return nil
}
