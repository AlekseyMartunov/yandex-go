// Package handlers contains all app s handlers
package handlers

import "net/http"

// DataBaseStatus uses to check db
func (s *ShortURLHandler) DataBaseStatus(w http.ResponseWriter, r *http.Request) {
	err := s.encoder.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
