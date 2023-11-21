// Package handlers contains all app s handlers
package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
)

type Str []string

// DeleteURL uses to delete url
func (s *ShortURLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserId")
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read request body error", http.StatusInternalServerError)
	}

	var data Str

	json.Unmarshal(b, &data)

	for _, val := range data {
		msg := simpleurl.URLToDel{
			UserID: userID,
			URL:    val,
		}
		s.delCh <- msg
	}

	w.WriteHeader(http.StatusAccepted)
}

// asyncDelURL uses to async deleting url
func (s *ShortURLHandler) asyncDelURL() {
	ticker := time.NewTicker(5 * time.Second)

	messages := make([]simpleurl.URLToDel, 0, 50)

	for {
		select {
		case msg, ok := <-s.delCh:
			if !ok {
				s.delCh = nil
				continue
			}
			messages = append(messages, msg)

		case <-ticker.C:
			if len(messages) == 0 {
				continue
			}
			err := s.encoder.DeleteURL(messages...)
			if err != nil {
				log.Fatalln(err)
				continue
			}
			messages = messages[:0]

		case <-s.ctx.Done():
			close(s.delCh)
			msg := <-s.delCh
			err := s.encoder.DeleteURL(msg)
			if err != nil {
				log.Fatalln(err)
			}
			return
		}
	}
}
