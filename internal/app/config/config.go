package config

import (
	"flag"
	"os"
)

type Config struct {
	addr            string `env:"SERVER_ADDRESS"`
	baseHost        string `env:"BASE_URL"`
	fileStoragePath string `env:"FILE_STORAGE_PATH"`
	dataBaseDSN     string `env:"DATABASE_DSN"`
	dataBaseStatus  bool
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

	flag.StringVar(&c.dataBaseDSN, "d", "",
		"Параметры БД")

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

	if val, ok := os.LookupEnv("DATABASE_DSN"); ok {
		c.dataBaseDSN = val
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

func (c *Config) GetDataBaseDSN() string {
	return c.dataBaseDSN
}

func (c *Config) GetDataBaseStatus() bool {
	return c.dataBaseStatus
}

func (c *Config) SetDataBaseStatus(status bool) {
	c.dataBaseStatus = status
}
