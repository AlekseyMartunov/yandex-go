package postgres

import (
	"context"
	"database/sql"
)

type URLModel struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *URLModel {
	return &URLModel{
		db: db,
	}
}

func (m *URLModel) CreateTableURL() error {
	query := `CREATE TABLE IF NOT EXISTS url (
    					id serial PRIMARY KEY, 
    					shorted VARCHAR(20),
    					original TEXT
    					)`

	_, err := m.db.ExecContext(context.Background(), query)
	return err

}

func (m *URLModel) Save(key, val string) error {
	query := "INSERT INTO url (shorted, original) VALUES ($1, $2)"

	_, err := m.db.ExecContext(context.Background(), query, key, val)
	if err != nil {
		return err
	}
	return nil
}

func (m *URLModel) Get(key string) (string, bool) {
	query := "SELECT original FROM url WHERE shorted = $1"
	row := m.db.QueryRowContext(context.Background(), query, key)

	var original sql.NullString
	row.Scan(&original)

	return original.String, original.Valid

}

func (m *URLModel) SaveBatch(data *[][3]string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	ctx := context.Background()
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO url (shorted, original) VALUES ($1, $2)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, val := range *data {
		_, err := stmt.ExecContext(ctx, val[2], val[1])
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (m *URLModel) Close() error {
	return m.db.Close()
}
