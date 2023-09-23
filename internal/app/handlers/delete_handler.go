package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
)

type Str []string

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

func (s *ShortURLHandler) asyncDelURL() {
	ticker := time.NewTicker(5 * time.Second)

	var messages []simpleurl.URLToDel

	for {
		select {
		case msg := <-s.delCh:
			messages = append(messages, msg)

		case <-ticker.C:
			if len(messages) == 0 {
				continue
			}
			err := s.encoder.DeleteURL(messages...)
			if err != nil {
				fmt.Println(err)
				continue
			}
			messages = nil
		}
	}
}
