package app

import (
	"flag"
	"os"
)

type netAdress struct {
	addr     string `env:"SERVER_ADDRESS"`
	baseHost string `env:"BASE_URL"`
}

func (a *app) GetConfig() {

	flag.StringVar(&a.cfg.addr, "a", "127.0.0.1:8080", "Адрес для запуска приложения")
	flag.StringVar(&a.cfg.baseHost, "b", "http://127.0.0.1:8080", "Базовый адрес сокращенного URL")

	flag.Parse()

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		a.cfg.addr = val
	}

	if val, ok := os.LookupEnv("BASE_URL"); ok {
		a.cfg.baseHost = val
	}
}

func (a *app) GetAdres() string {
	return a.cfg.addr
}

func (a *app) GetShorterURL() string {
	return a.cfg.baseHost + "/"
}
