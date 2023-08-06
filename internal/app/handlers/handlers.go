package handlers

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (s *Server) EncodeURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	data, err := io.ReadAll(r.Body)
	if err != nil || string(data) == "" {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}
	id := s.Cfg.BaseAddr + "/" + s.db.Encode(string(data))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (s *Server) DecodeURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "url_id")
	url, ok := s.db.Decode(id)

	if !ok {
		http.Error(w, "Empty key", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
