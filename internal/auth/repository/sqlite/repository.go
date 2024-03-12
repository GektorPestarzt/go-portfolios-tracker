package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"go-portfolios-tracker/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Get(ctx context.Context, username, password string) (*models.User, error) {
	row := ur.db.QueryRow(`SELECT uuid, username, password FROM portfolios WHERE username = $1`, username)

	user := &models.User{}
	err := row.Scan(&user.UUID, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	if password != user.Password {
		// TODO: refactor error creating
		return nil, fmt.Errorf("incorrect passport")
	}

	return user, nil
}

func (ur *UserRepository) Add(ctx context.Context, user *models.User) error {
	_, err := ur.db.Exec(`INSERT INTO portfolios (username, password) VALUES (?, ?)`,
		user.Username,
		user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) Delete(ctx context.Context, uuid int) error {
	return nil
}
