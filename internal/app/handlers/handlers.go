// Package handlers contains all app s handlers
package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"

	"github.com/go-chi/chi/v5"
)

type encoder interface {
	Encode(url, userID string) (string, error)
	Decode(string) (string, error)
	BatchEncode(data *[][3]string, userID string) error
	GetShorted(key string) (string, bool)
	GetAllURL(userID string) ([][2]string, error)
	DeleteURL(...simpleurl.URLToDel) error
	Ping() error
}

type config interface {
	GetShorterURL() string
	GetDataBaseStatus() bool
}

// ShortURLHandler type store info about handles app
type ShortURLHandler struct {
	encoder encoder
	cfg     config
	ctx     context.Context
	delCh   chan simpleurl.URLToDel
}

// NewShortURLHandler return new ShortURLHandler
func NewShortURLHandler(e encoder, c config, ctx context.Context) *ShortURLHandler {
	h := ShortURLHandler{
		encoder: e,
		cfg:     c,
		ctx:     ctx,
	}
	h.delCh = make(chan simpleurl.URLToDel, 1024)

	go h.asyncDelURL()
	return &h
}

// EncodeURL encode url
func (s *ShortURLHandler) EncodeURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	data, err := io.ReadAll(r.Body)
	if err != nil || string(data) == "" {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("userID")

	encodedData, err := s.encoder.Encode(string(data), userID)
	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
		shorted, ok := s.encoder.GetShorted(string(data))
		if ok {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(s.cfg.GetShorterURL() + shorted))
			return
		}
		http.Error(w, "Some server error", http.StatusInternalServerError)
		return
	}

	id := s.cfg.GetShorterURL() + encodedData
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

// DecodeURL decode url handler
func (s *ShortURLHandler) DecodeURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "url_id")
	url, err := s.encoder.Decode(id)

	if err != nil {
		if errors.Is(err, simpleurl.ErrDeletedURL) {
			http.Error(w, "Deleted key ", http.StatusGone)
			return

		}
		if errors.Is(err, simpleurl.ErrEmptyKey) {
			http.Error(w, "Empty key", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
