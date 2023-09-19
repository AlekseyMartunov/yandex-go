package handlers

import (
	"encoding/json"
	"net/http"
)

type allURl struct {
	Shorted  string `json:"short_url"`
	Original string `json:"original_url"`
}

type arrAllURL []allURl

func (s *ShortURLHandler) GetAllURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")
	data, err := s.encoder.GetAllURL(userID)
	if err != nil {
		http.Error(w, "Some server error", http.StatusInternalServerError)
	}

	if len(data) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result := make(arrAllURL, 0, len(data))

	for indx := range data {
		el := allURl{
			Shorted:  data[indx][0],
			Original: data[indx][1],
		}
		result = append(result, el)
	}

	marshalRes, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Some server error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(marshalRes)
}