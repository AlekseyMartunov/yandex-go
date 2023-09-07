package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type encoder interface {
	Encode(string) string
	Decode(string) (string, bool)
	BatchEncode(*[][3]string) error
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

	id := s.cfg.GetShorterURL() + s.encoder.Encode(string(data))
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
