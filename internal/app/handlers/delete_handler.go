package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type Str []string

func (s *ShortURLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read request body error", http.StatusInternalServerError)
	}

	var data Str
	json.Unmarshal(b, &data)

	ch := fanOut(data)
	err = s.encoder.DeleteURLByUserID(userID, r.Context(), ch)
	if err != nil {
		http.Error(w, "some error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func fanOut(data []string) chan string {
	nCh := len(data)
	var wg sync.WaitGroup

	ch := make(chan string, nCh)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, val := range data {
			ch <- val
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
