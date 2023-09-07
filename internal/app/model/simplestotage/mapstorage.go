package simplestotage

import "sync"

type MapStorage struct {
	data map[string]string
	sync.Mutex
}

func (s *MapStorage) Save(key, val string) error {
	s.Mutex.Lock()
	s.data[key] = val
	s.Mutex.Unlock()
	return nil
}

func (s *MapStorage) Get(key string) (string, bool) {
	val, ok := s.data[key]
	return val, ok
}

func (s *MapStorage) SaveBatch(data *[][3]string) error {
	// [[a, b, c], [a, b, c], ...]
	// a - CorrelationID
	// b - OriginalURL
	// c - ShortedURL

	for id := range *data {
		key := (*data)[id][2]
		val := (*data)[id][1]
		s.Save(key, val)
	}
	return nil
}

func (s *MapStorage) Close() error {
	return nil
}
