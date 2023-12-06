package encoder

import (
	"fmt"
	"log"

	"github.com/AlekseyMartunov/yandex-go.git/internal/app/model/url/simpleurl"
)

func Example() {
	url1 := "some string here..."
	url2 := "another string here...."

	userID := "1"

	storage, err := simpleurl.NewMapStorage()
	if err != nil {
		log.Fatalln(err)
	}

	encoder := NewEncoder(storage)

	short1, err := encoder.Encode(url1, userID)
	if err != nil {
		log.Fatalln(err)
	}

	short2, err := encoder.Encode(url2, userID)
	if err != nil {
		log.Fatalln(err)
	}

	origing1, err := encoder.Decode(short1)
	if err != nil {
		log.Fatalln(err)
	}

	origing2, err := encoder.Decode(short2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(origing1)
	fmt.Println(origing2)

	//Output:
	//some string here...
	//another string here....
}
