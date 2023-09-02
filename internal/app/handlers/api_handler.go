package handlers

import (
	"encoding/json"
	"io"
	"net/http"
)

type jsonRequest struct {
	URL string `json:"url"`
}

type jsonResponse struct {
	Result string `json:"result"`
}

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

	jResp := jsonResponse{Result: s.cfg.GetShorterURL() + s.encoder.Encode(jReq.URL)}

	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(jResp)

	if err != nil {
		panic(err)
	}
	w.Write(res)

}
