package config

import "fmt"

func Example() {
	cfg := NewConfig()
	cfg.GetConfig()

	fmt.Println(cfg.GetAddress())
	fmt.Println(cfg.GetShorterURL())

	//Output:
	//127.0.0.1:8080
	//http://127.0.0.1:8080/
}
