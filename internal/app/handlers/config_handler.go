package handlers

import "net/http"

func (s *ShortURLHandler) DataBaseStatus(w http.ResponseWriter, r *http.Request) {
	if s.cfg.GetDataBaseStatus() {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
