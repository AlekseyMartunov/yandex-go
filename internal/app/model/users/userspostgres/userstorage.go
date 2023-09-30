package userspostgres

import (
	"context"
	"database/sql"
)

type UserModel struct {
	db *sql.DB
}

func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

func (u *UserModel) GetFreeID() (int, error) {
	query := `SELECT COUNT(*) FROM users`
	row := u.db.QueryRowContext(context.Background(), query)
	var count int

	err := row.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (u *UserModel) SaveNewUser() (int, error) {
	query := `INSERT INTO users VALUES (DEFAULT)`
	_, err := u.db.ExecContext(context.Background(), query)
	if err != nil {
		return -1, err
	}
	return -1, nil
}
