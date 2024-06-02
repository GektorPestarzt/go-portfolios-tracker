package sqlite

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
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

type portfolioBytes struct {
	ID                    int
	TotalAmountPortfolio  []byte
	TotalAmountShares     []byte
	TotalAmountBonds      []byte
	TotalAmountEtf        []byte
	TotalAmountCurrencies []byte
	ExpectedYield         []byte
}

type positionBytes struct {
	ID                   int
	Figi                 string
	InstrumentType       string
	Quantity             []byte
	AveragePositionPrice []byte
	ExpectedYield        []byte
}

// TODO: remake. Test version
func (pr *PortfolioRepository) Init(ctx context.Context, username string, token string, typeA string) (int64, error) {
	res, err := pr.db.Exec(`INSERT INTO accounts (token, type, username, status) VALUES (?, ?, ?, ?)`,
		token,
		typeA,
		username,
		models.Created)
	// TODO: remake types of accounts

	if err != nil {
		// TODO: HA
		return -1, err
	}

	return res.LastInsertId()
}

func (pr *PortfolioRepository) Get(ctx context.Context, id int64) (*models.Account, error) {
	row := pr.db.QueryRow(`SELECT * FROM accounts WHERE id = ?`, id)
	account := &models.Account{}
	err := row.Scan(&account.ID, &account.Token, &account.Type, &account.Username, &account.Status)
	if err != nil {
		// TODO: HA
		return nil, err
	}

	rows, err := pr.db.Query(`SELECT * FROM portfolios WHERE account_id = ?`, id)
	if err != nil {
		// TODO: HA
		return nil, err
	}

	for rows.Next() {
		var portfolioID, aID int
		portfolioB := &portfolioBytes{}
		if err := rows.Scan(&portfolioID, &portfolioB.TotalAmountPortfolio,
			&portfolioB.TotalAmountShares, &portfolioB.TotalAmountBonds,
			&portfolioB.TotalAmountEtf, &portfolioB.TotalAmountCurrencies,
			&portfolioB.ExpectedYield, &aID); err != nil {
			// TODO
			return nil, err
		}
		var portfolio models.Portfolio
		err := portfolioConvertFrom(&portfolio, portfolioB)
		if err != nil {
			// TODO
			return nil, err
		}

		rowsPos, err := pr.db.Query(`SELECT * FROM positions WHERE portfolio_id = ?`, portfolioID)
		if err != nil {
			// TODO: HA
			return nil, err
		}

		for rowsPos.Next() {
			var positionID, pID int
			positionB := &positionBytes{}
			if err := rowsPos.Scan(&positionID, &positionB.Figi, &positionB.InstrumentType,
				&positionB.Quantity, &positionB.AveragePositionPrice,
				&positionB.ExpectedYield, &pID); err != nil {
				// TODO
				return nil, err
			}
			var position models.Position
			err := positionConvertFrom(&position, positionB)
			if err != nil {
				// TODO
				return nil, err
			}

			portfolio.Positions = append(portfolio.Positions, &position)
		}

		account.Portfolios = append(account.Portfolios, &portfolio)
	}

	return account, nil
}

func (pr *PortfolioRepository) Update(ctx context.Context, account *models.Account) error {
	for _, portfolio := range account.Portfolios {
		portfolioB, err := portfolioConvertTo(portfolio)
		if err != nil {
			// TODO
			return err
		}

		res, err := pr.db.Exec(`INSERT INTO portfolios (
        total_amount_portfolio,
		total_amount_shares,
		total_amount_bonds,
		total_amount_etf,
		total_amount_currencies,
		expected_yield,
		account_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			portfolioB.TotalAmountPortfolio,
			portfolioB.TotalAmountShares,
			portfolioB.TotalAmountBonds,
			portfolioB.TotalAmountEtf,
			portfolioB.TotalAmountCurrencies,
			portfolioB.ExpectedYield,
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
			positionB, err := positionConvertTo(position)
			if err != nil {
				// TODO
				return err
			}

			_, err = pr.db.Exec(`INSERT INTO positions (
            figi,
            instrument_type,
            quantity,
			average_position_price,
			expected_yield,
			portfolio_id) VALUES (?, ?, ?, ?, ?, ?)`,
				positionB.Figi,
				positionB.InstrumentType,
				positionB.Quantity,
				positionB.AveragePositionPrice,
				positionB.ExpectedYield,
				id)

			if err != nil {
				// TODO
				return err
			}
		}
	}

	return nil
}

func (pr *PortfolioRepository) UpdateStatus(ctx context.Context, id int64, status models.Status) error {
	_, err := pr.db.Exec(`UPDATE accounts SET status = ? WHERE id = ?`,
		status,
		id)

	if err != nil {
		return err
	}
	return nil
}

func (pr *PortfolioRepository) GetToken(ctx context.Context, id int64) (string, error) {
	return "", nil
}

func (pr *PortfolioRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func portfolioConvertTo(src *models.Portfolio) (*portfolioBytes, error) {
	var pb portfolioBytes
	var err error

	pb.TotalAmountPortfolio = encodeToBytes(src.TotalAmountPortfolio, &err)
	pb.TotalAmountShares = encodeToBytes(src.TotalAmountShares, &err)
	pb.TotalAmountBonds = encodeToBytes(src.TotalAmountBonds, &err)
	pb.TotalAmountEtf = encodeToBytes(src.TotalAmountEtf, &err)
	pb.TotalAmountCurrencies = encodeToBytes(src.TotalAmountCurrencies, &err)
	pb.ExpectedYield = encodeToBytes(src.ExpectedYield, &err)

	if err != nil {
		return nil, err
	}
	return &pb, nil
}

func portfolioConvertFrom(dst *models.Portfolio, src *portfolioBytes) error {
	var err error

	dst.TotalAmountPortfolio = decodeToMoneyValue(src.TotalAmountPortfolio, &err)
	dst.TotalAmountShares = decodeToMoneyValue(src.TotalAmountShares, &err)
	dst.TotalAmountBonds = decodeToMoneyValue(src.TotalAmountBonds, &err)
	dst.TotalAmountEtf = decodeToMoneyValue(src.TotalAmountEtf, &err)
	dst.TotalAmountCurrencies = decodeToMoneyValue(src.TotalAmountCurrencies, &err)
	dst.ExpectedYield = decodeToQuotation(src.ExpectedYield, &err)

	if err != nil {
		return err
	}
	return nil
}

func positionConvertTo(src *models.Position) (*positionBytes, error) {
	var pb positionBytes
	var err error

	pb.Figi = src.Figi
	pb.InstrumentType = src.InstrumentType
	pb.Quantity = encodeToBytes(src.Quantity, &err)
	pb.AveragePositionPrice = encodeToBytes(src.AveragePositionPrice, &err)
	pb.ExpectedYield = encodeToBytes(src.ExpectedYield, &err)

	if err != nil {
		return nil, err
	}
	return &pb, nil
}

func positionConvertFrom(dst *models.Position, src *positionBytes) error {
	var err error

	dst.Figi = src.Figi
	dst.InstrumentType = src.InstrumentType
	dst.Quantity = decodeToQuotation(src.Quantity, &err)
	dst.AveragePositionPrice = decodeToMoneyValue(src.AveragePositionPrice, &err)
	dst.ExpectedYield = decodeToQuotation(src.ExpectedYield, &err)

	if err != nil {
		return err
	}
	return nil
}

func encodeToBytes(p interface{}, err *error) []byte {
	if *err != nil {
		return nil
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	*err = enc.Encode(p)

	return buf.Bytes()
}

func decodeToMoneyValue(src []byte, err *error) *models.MoneyValue {
	if *err != nil {
		return nil
	}

	var mv models.MoneyValue
	dec := gob.NewDecoder(bytes.NewReader(src))
	*err = dec.Decode(&mv)

	return &mv
}

func decodeToQuotation(src []byte, err *error) *models.Quotation {
	if *err != nil {
		return nil
	}

	var q models.Quotation
	dec := gob.NewDecoder(bytes.NewReader(src))
	*err = dec.Decode(&q)

	return &q
}

/*
func decodeFromBytes(src []byte, err *error) interface{} {
	if *err != nil {
		return nil
	}

	var p interface{}
	dec := gob.NewDecoder(bytes.NewReader(src))
	*err = dec.Decode(&p)

	return p
}
*/
