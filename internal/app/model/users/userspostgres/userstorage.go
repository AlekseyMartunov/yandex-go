// Package userspostgres store struct to works whit postgres connection and query
package userspostgres

import (
	"context"

	"database/sql"
)

// UserModel store postgres connections
type UserModel struct {
	db *sql.DB
}

// NewUserModel create new user instance
func NewUserModel(db *sql.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

// SaveNewUser save new user and return id
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
