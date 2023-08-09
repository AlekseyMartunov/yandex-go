package config

import (
	"flag"
	"os"
)

type config struct {
	addr     string `env:"SERVER_ADDRESS"`
	baseHost string `env:"BASE_URL"`
}

func NewConfig() *config {
	c := config{}
	return &c
}

func (c *config) GetConfig() {

	flag.StringVar(&c.addr, "a", "127.0.0.1:8080", "Адрес для запуска приложения")
	flag.StringVar(&c.baseHost, "b", "http://127.0.0.1:8080", "Базовый адрес сокращенного URL")

	flag.Parse()

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.addr = val
	}

	if val, ok := os.LookupEnv("BASE_URL"); ok {
		c.baseHost = val
	}
}

func (c *config) GetAddress() string {
	return c.addr
}

func (c *config) GetShorterURL() string {
	return c.baseHost + "/"
}
