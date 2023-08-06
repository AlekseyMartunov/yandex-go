package app

import "flag"

type netAdress struct {
	addr     string
	baseHost string
}

func (a *app) GetConfig() {
	flag.StringVar(&a.cfg.addr, "a", "127.0.0.1:8080", "Адрес для запуска приложения")
	flag.StringVar(&a.cfg.baseHost, "b", "http://localhost:8080/", "Базовый адрес сокращенного URL")

	flag.Parse()
}

func (a *app) GetAdres() string {
	return a.cfg.addr
}

func (a *app) GetShorterURL() string {
	return a.cfg.baseHost
}
