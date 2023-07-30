package app

import (
	"io"
	"net/http"
)

func (s *app) ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.URL.String() != "/" {
			http.Error(w, "You should send a request to '/'", http.StatusBadRequest)
			return
		}
		data, _ := io.ReadAll(r.Body)
		id := s.encode(string(data))
		w.Header().Set("content-type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(id))
	}

	if r.Method == http.MethodGet {
		id := r.URL.String()[1:]
		url, ok := s.decode(id)
		if !ok {
			http.Error(w, "Empty key", http.StatusBadRequest)
			return
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}

}
