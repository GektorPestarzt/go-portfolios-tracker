package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"go-portfolios-tracker/internal/models"
	"log"
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
func (pr *PortfolioRepository) Init(ctx context.Context, username string, token string, typeA string) error {
	_, err := pr.db.Exec(`INSERT INTO accounts (token, type, username) VALUES (?, ?, ?)`,
		token,
		typeA,
		username)
	// TODO: remake types of accounts

	if err != nil {
		// TODO: HA
		return err
	}

	return nil
}

func (pr *PortfolioRepository) Get(ctx context.Context, id int) (*models.Account, error) {
	row := pr.db.QueryRow(`SELECT * FROM accounts WHERE id = ?`, id)
	account := &models.Account{}
	err := row.Scan(&account.ID, &account.Token, &account.Type, &account.Username)

	if err != nil {
		// TODO: HA
		return nil, err
	}

	return account, nil
}

func (pr *PortfolioRepository) Update(ctx context.Context, account *models.Account) error {
	for _, portfolio := range account.Portfolios {
		res, err := pr.db.Exec(`INSERT INTO portfolios (
        total_amount_portfolio,
		total_amount_shares,
		total_amount_bonds,
		total_amount_etf,
		total_amount_currencies,
		expected_yield,
		account_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			EncodeToBytes(portfolio.TotalAmountPortfolio),
			EncodeToBytes(portfolio.TotalAmountShares),
			EncodeToBytes(portfolio.TotalAmountBonds),
			EncodeToBytes(portfolio.TotalAmountEtf),
			EncodeToBytes(portfolio.TotalAmountCurrencies),
			EncodeToBytes(portfolio.ExpectedYield),
			account.ID)

		if err != nil {
			// TODO
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			// TODO
			return err
		}

		for _, position := range portfolio.Positions {
			_, err := pr.db.Exec(`INSERT INTO positions (
            figi,
            instrument_type,
            quantity,
			average_position_price,
			expected_yield,
			portfolio_id) VALUES (?, ?, ?, ?, ?, ?)`,
				position.Figi,
				position.InstrumentType,
				EncodeToBytes(position.Quantity),
				EncodeToBytes(position.AveragePositionPrice),
				EncodeToBytes(position.ExpectedYield),
				id)

			if err != nil {
				// TODO
				return err
			}
		}
	}

	return nil
}

func (pr *PortfolioRepository) GetToken(ctx context.Context, id int) (string, error) {
	return "", nil
}

func (pr *PortfolioRepository) Delete(ctx context.Context, id int) error {
	return nil
}

func EncodeToBytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	err := enc.Encode(p)
	if err != nil {
		// TODO
		log.Fatal(err)
	}

	return buf.Bytes()
}
