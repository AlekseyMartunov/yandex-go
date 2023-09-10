package handlers

import (
	"github.com/jackc/pgx/v5/pgconn"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type encoder interface {
	Encode(string) (string, error)
	Decode(string) (string, bool)
	BatchEncode(*[][3]string) error
	GetShorted(key string) (string, bool)
}

type config interface {
	GetShorterURL() string
	GetDataBaseStatus() bool
}

type ShortURLHandler struct {
	encoder encoder
	cfg     config
}

func NewShortURLHandler(e encoder, c config) *ShortURLHandler {
	return &ShortURLHandler{encoder: e, cfg: c}
}

func (s *ShortURLHandler) EncodeURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	data, err := io.ReadAll(r.Body)
	if err != nil || string(data) == "" {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	encodedData, err := s.encoder.Encode(string(data))
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

func (s *ShortURLHandler) DecodeURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "url_id")
	url, ok := s.encoder.Decode(id)

	if !ok {
		http.Error(w, "Empty key", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
