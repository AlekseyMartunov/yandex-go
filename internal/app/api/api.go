package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type encoder interface {
	Encode(string) string
	Decode(string) (string, bool)
}

type config interface {
	GetShorterURL() string
}

type api struct {
	encoder encoder
	cfg     config
}

type jsonRequest struct {
	URL string `json:"url"`
}

type jsonResponse struct {
	Result string `json:"result"`
}

func NewAPIHandlers(e encoder, c config) *api {
	return &api{encoder: e, cfg: c}
}

func (a *api) EncodeAPI(w http.ResponseWriter, r *http.Request) {
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

	JResp := jsonResponse{Result: a.cfg.GetShorterURL() + a.encoder.Encode(jReq.URL)}

	w.WriteHeader(http.StatusCreated)

	res, err := json.Marshal(JResp)

	if err != nil {
		panic(err)
	}
	w.Write(res)

}
