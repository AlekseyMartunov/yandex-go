package encoder

import (
	storage2 "github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
	"testing"
)

func BenchmarkEncoder_Decode(b *testing.B) {

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

	b.ResetTimer()
	b.StopTimer()

	for id := range testCase {
		b.StartTimer()

		key := keys[id]
		e.Decode(key)

		b.StopTimer()
	}
}
