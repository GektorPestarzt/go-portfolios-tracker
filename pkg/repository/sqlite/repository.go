package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(repositoryPath string) (*sql.DB /**Repository*/, error) {
	const op = "repository.sqlite.New"

	db, err := sql.Open("sqlite3", repositoryPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	/*
		stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users(
		    uuid INTEGER PRIMARY KEY,
		    username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS username_index ON portfolios(username);
		`)
	*/

	/*
		stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS accounts (
			id INTEGER PRIMARY KEY,
			token TEXT,
			type TEXT NOT NULL,
			portfolios_id INTEGER,
			username TEXT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(uuid)
		);
		CREATE TABLE IF NOT EXISTS portfolios (
			id INTEGER PRIMARY KEY,
			total_amount_portfolio BLOB,
			total_amount_shares BLOB,
			total_amount_bonds BLOB,
			total_amount_etf BLOB,
			total_amount_currencies BLOB,
			expected_yield BLOB,
			positions_id INTEGER,
			FOREIGN KEY (id) REFERENCES accounts(portfolios_id)
		);
		CREATE TABLE IF NOT EXISTS positions (
			id INTEGER PRIMARY KEY,
			figi TEXT,
			instrument_type TEXT,
			quantity BLOB,
			average_position_price BLOB,
			expected_yield BLOB,
			FOREIGN KEY (id) REFERENCES portfolios(positions_id)
		);
		`)

		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}*/

	// _, err = stmt.Exec()
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY,
		token TEXT,
		type TEXT NOT NULL,
		username TEXT NOT NULL,
		FOREIGN KEY (username) REFERENCES users(username)
	);
	CREATE TABLE IF NOT EXISTS portfolios (
		id INTEGER PRIMARY KEY,
		total_amount_portfolio BLOB,
		total_amount_shares BLOB,
		total_amount_bonds BLOB,
		total_amount_etf BLOB,
		total_amount_currencies BLOB,
		expected_yield BLOB,
		account_id INTEGER,
		FOREIGN KEY (account_id) REFERENCES accounts(id)
	);
	CREATE TABLE IF NOT EXISTS positions (
		id INTEGER PRIMARY KEY,
		figi TEXT,
		instrument_type TEXT,
		quantity BLOB,
		average_position_price BLOB,
		expected_yield BLOB,
		portfolio_id INTEGER,
		FOREIGN KEY (portfolio_id) REFERENCES portfolios(id)
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	// return &Repository{db: db}, nil
	return db, nil
}
