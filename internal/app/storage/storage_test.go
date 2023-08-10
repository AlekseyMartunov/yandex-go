package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage(t *testing.T) {

	s := NewStorage()

	testCase := []string{
		"AAA",
		"BBB",
		"CCC",
		"474848",
		"https://www.google.com",
		"шла Cаша по шоссе",
		"Some another string",
		"\n\n\t\\dfdft\t\nkdfkdf",
		"AAA",
		"123",
	}

	keys := make([]string, 0, 10)

	for _, val := range testCase {
		key := s.Encode(val)
		keys = append(keys, key)
	}

	for id, val := range testCase {
		key := keys[id]
		res, ok := s.Decode(key)

		assert.True(t, ok, fmt.Sprintf("Ключ %s, не найден в мапе!", val))

		assert.Equal(t, val, res, fmt.Sprintf("Мапа вернула: %s, а нужно: %s", res, val))

	}

}
