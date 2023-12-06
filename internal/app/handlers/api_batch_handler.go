// Package handlers contains all app s handlers
package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

// batchRequest type uses to unmarshal http body
type batchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// batchResponse type uses to unmarshal http body
type batchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type jsonBatchReq []batchRequest
type jsonBatchResp []batchResponse

// BatchURL shorts any urls at ones
func (s *ShortURLHandler) BatchURL(w http.ResponseWriter, r *http.Request) {
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

	var jbReq jsonBatchReq

	err = json.Unmarshal(data, &jbReq)
	if err != nil {
		http.Error(w, "Request body read error", http.StatusBadRequest)
		return
	}

	urlArr := make([][3]string, 0, len(jbReq))

	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for _, val := range jbReq {
		if val.CorrelationID == "" || val.OriginalURL == "" {
			continue
		}
		urlArr = append(urlArr, [3]string{val.CorrelationID, val.OriginalURL, ""})
	}

	userID := r.Header.Get("userID")

	err = s.encoder.BatchEncode(&urlArr, userID)
	if err != nil {
		http.Error(w, "Saving error", http.StatusBadRequest)
		return
	}

	jbResp := make(jsonBatchResp, len(urlArr))

	for id, val := range urlArr {
		jbResp[id].ShortURL = s.cfg.GetShorterURL() + val[2]
		jbResp[id].CorrelationID = val[0]
	}

	res, err := json.Marshal(jbResp)

	if err != nil {
		http.Error(w, "Request body read error", http.StatusHTTPVersionNotSupported)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
