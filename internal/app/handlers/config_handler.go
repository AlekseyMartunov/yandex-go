package handlers

import "net/http"

func (s *ShortURLHandler) DataBaseStatus(w http.ResponseWriter, r *http.Request) {
	err := s.encoder.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
