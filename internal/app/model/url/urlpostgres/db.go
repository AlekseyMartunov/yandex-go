package urlpostgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
)

type URLModel struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *URLModel {
	return &URLModel{
		db: db,
	}
}

func (m *URLModel) Ping() error {
	return m.db.Ping()
}

func (m *URLModel) Save(key, val, userID string) error {
	query := "INSERT INTO url (shorted, original, user_id) VALUES ($1, $2, $3)"
	_, err := m.db.ExecContext(context.Background(), query, key, val, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return pgErr
		}
		return err
	}
	return nil
}

func (m *URLModel) Get(key string) (string, error) {
	query := "SELECT original, deleted FROM url WHERE shorted = $1"
	row := m.db.QueryRowContext(context.Background(), query, key)

	var original sql.NullString
	var flag bool
	row.Scan(&original, &flag)

	if flag {
		return "", simpleurl.ErrDeletedURL
	}

	if original.Valid {
		return original.String, nil
	}

	return original.String, simpleurl.ErrEmptyKey
}

func (m *URLModel) SaveBatch(data *[][3]string, userID string) error {
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

	query := `INSERT INTO url (shorted, original, user_id) VALUES ($1, $2, $3)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, val := range *data {
		_, err := stmt.ExecContext(ctx, val[2], val[1], userID)
		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (m *URLModel) GetShorted(key string) (string, bool) {
	query := "SELECT shorted FROM url WHERE original = $1"
	row := m.db.QueryRowContext(context.Background(), query, key)

	var original sql.NullString
	row.Scan(&original)

	return original.String, original.Valid
}

func (m *URLModel) GetAllURL(userID string) ([][2]string, error) {
	query := `SELECT shorted, original FROM url WHERE user_id = $1`
	rows, err := m.db.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result [][2]string

	for rows.Next() {
		var short string
		var origin string

		err = rows.Scan(&short, &origin)
		if err != nil {
			return nil, err
		}

		var arr = [2]string{short, origin}
		result = append(result, arr)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (m *URLModel) DeleteURL(messages ...simpleurl.URLToDel) error {
	ctx := context.Background()

	query := `UPDATE url 
				SET deleted = TRUE
				WHERE shorted = $1 AND user_id = $2;`

	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, msg := range messages {
		_, err := stmt.ExecContext(ctx, msg.URL, msg.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *URLModel) Close() error {
	return m.db.Close()
}
