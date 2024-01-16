package handlers

import (
	"encoding/json"
	"net"
	"net/http"
)

// statResponse support struct
type statResponse struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// Stats show how many users and url contains application
func (s *ShortURLHandler) Stats(w http.ResponseWriter, r *http.Request) {
	IP := s.cfg.GetTrustedIP()
	if IP == nil {
		http.Error(w, "I do not know your ip address", http.StatusForbidden)
	}

	userIP := net.ParseIP(r.Header.Get("X-Real-IP"))
	if !IP.Contains(userIP) {
		http.Error(w, "I do not know your ip address", http.StatusForbidden)
	}

	urls, users := s.encoder.Statistics()
	resp := statResponse{
		Urls:  urls,
		Users: users,
	}

	res, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "some error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
