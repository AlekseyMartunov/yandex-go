package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type storage interface {
	Encode(url string) string
	Decode(shortURL string) (string, bool)
}

type config interface {
	GetShorterURL() string
}

type ShortURLHandler struct {
	storage storage
	cfg     config
}

func NewShortURLHandler(storage storage, cfg config) *ShortURLHandler {
	return &ShortURLHandler{storage: storage, cfg: cfg}
}

func (s *ShortURLHandler) EncodeURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	data, err := io.ReadAll(r.Body)
	if err != nil || string(data) == "" {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	id := s.cfg.GetShorterURL() + s.storage.Encode(string(data))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *ShortURLHandler) DecodeURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "url_id")
	url, ok := s.storage.Decode(id)

	if !ok {
		http.Error(w, "Empty key", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
