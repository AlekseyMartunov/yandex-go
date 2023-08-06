package config

import "flag"

type NetAddress struct {
	Host     string
	BaseAddr string
}

func Config() *NetAddress {
	add := &NetAddress{}
	flag.StringVar(&add.Host, "a", "127.0.0.1:8080", "Хост и порт для запуска приложения")
	flag.StringVar(&add.BaseAddr, "b", "localhost:8080", "базовый адрес результирующего сокращённого URL ")

	flag.Parse()

	return add
}
