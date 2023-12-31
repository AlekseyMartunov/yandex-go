package encoder

import (
	"fmt"
	"testing"

	storage2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {

	s, _ := storage2.NewMapStorage()

	e := NewEncoder(s)

	testCase := []string{
		"AAA",
		"BBB",
		"CCC",
		"474848",
		"https://www.google.com",
		"шла Cаша по шоссе",
		"Some another string",
		"\n\n\t\\dfdft\t\nkdfkdf",
		"123",
	}

	keys := make([]string, 0, 10)

	for _, val := range testCase {
		key, _ := e.Encode(val, "10")
		keys = append(keys, key)
	}

	for id, val := range testCase {
		key := keys[id]
		res, err := e.Decode(key)

		assert.Equal(t, nil, err, fmt.Sprintf("Ключ %s, не найден в мапе!", val))

		assert.Equal(t, val, res, fmt.Sprintf("Мапа вернула: %s, а нужно: %s", res, val))

	}

}
