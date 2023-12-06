// Package handlers contains all app s handlers
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// jsonRequest type uses to unmarshal http body
type jsonRequest struct {
	URL string `json:"url"`
}

// jsonResponse type uses to unmarshal http body
type jsonResponse struct {
	Result string `json:"result"`
}

// EncodeAPI encode url
func (s *ShortURLHandler) EncodeAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-type should be 'application/json'", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Request body read error", http.StatusBadRequest)
		return
	}

	jReq := jsonRequest{}
	err = json.Unmarshal(data, &jReq)
	if err != nil {
		http.Error(w, "Request body read error", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("userID")

	encodedData, err := s.encoder.Encode(jReq.URL, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			shorted, ok := s.encoder.GetShorted(jReq.URL)
			if ok {
				jResp := jsonResponse{
					Result: s.cfg.GetShorterURL() + shorted,
				}

				res, err := json.Marshal(jResp)
				if err != nil {
					http.Error(w, "Request body read error", http.StatusHTTPVersionNotSupported)
					return
				}
				w.WriteHeader(http.StatusConflict)
				w.Write(res)
				return
			}
			http.Error(w, "Some server error", http.StatusInternalServerError)
			return
		}
	}

	jResp := jsonResponse{
		Result: s.cfg.GetShorterURL() + encodedData,
	}

	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(jResp)

	if err != nil {
		http.Error(w, "Request body read error", http.StatusHTTPVersionNotSupported)
		return
	}
	w.Write(res)
}
