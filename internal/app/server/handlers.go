package server

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func (br *baseRouter) encodeURL(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/plain")
	data, err := io.ReadAll(r.Body)
	if err != nil || string(data) == "" {
		http.Error(w, "Missing body", http.StatusBadRequest)
		return
	}

	id := br.cfg.GetShorterURL() + br.db.Encode(string(data))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}

func (br *baseRouter) decodeURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "url_id")
	url, ok := br.db.Decode(id)

	if !ok {
		http.Error(w, "Empty key", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
