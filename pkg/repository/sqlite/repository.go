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

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS portfolios(
	    uuid INTEGER PRIMARY KEY,
	    username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS username_index ON portfolios(username);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	// return &Repository{db: db}, nil
	return db, nil
}
