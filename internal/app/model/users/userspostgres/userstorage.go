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

func (u *UserModel) SaveNewUser() (int, error) {
	query := `INSERT INTO users VALUES (DEFAULT) RETURNING user_id`
	row := u.db.QueryRowContext(context.Background(), query)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return userID, err
	}
	return userID, nil
}
