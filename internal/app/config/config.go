package config

import (
	"flag"
	"os"
)

type Config struct {
	addr            string `env:"SERVER_ADDRESS"`
	baseHost        string `env:"BASE_URL"`
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GetConfig() {

	flag.StringVar(&c.addr, "a", "127.0.0.1:8080",
		"Адрес для запуска приложения")

	flag.StringVar(&c.baseHost, "b", "http://127.0.0.1:8080",
		"Базовый адрес сокращенного URL")

	flag.StringVar(&c.fileStoragePath, "f", "/tmp/short-url-db.json",
		"Путь до файла-хранилища")

	flag.Parse()

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.addr = val
	}

	if val, ok := os.LookupEnv("BASE_URL"); ok {
		c.baseHost = val
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.fileStoragePath = val
	}
}

func (c *Config) GetAddress() string {
	return c.addr
}

func (c *Config) GetShorterURL() string {
	return c.baseHost + "/"
}

func (c *Config) GetFileStoragePath() string {
	return c.fileStoragePath
}
