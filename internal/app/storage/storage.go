package storage

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{data: make(map[string]string)}
}

func (s *Storage) Save(key, val string) {
	s.data[key] = val
}

func (s *Storage) Get(key string) (string, bool) {
	res, ok := s.data[key]
	return res, ok
}
