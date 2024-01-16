// Package config  store information about flags application
package config

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"os"
	"strconv"
)

// Config type store information about flags application
type Config struct {
	Addr            string `env:"SERVER_ADDRESS" json:"server_address"`
	BaseHost        string `env:"BASE_URL" json:"base_url"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DataBaseDSN     string `env:"DATABASE_DSN" json:"database_dsn"`
	HTTPS           bool   `env:"ENABLE_HTTPS" json:"enable_https"`
	DataBaseStatus  bool
	TrustedIP       *net.IPNet
}

// supportIP help with convert json value to Config store
type supportIP struct {
	IP string `json:"trusted_subnet"`
}

// NewConfig create new config struct
func NewConfig() *Config {
	return &Config{}
}

// GetConfig update values of struct
func (c *Config) GetConfig() {
	var fileName string

	flag.StringVar(&fileName, "c", "", "фаил конфигурации")
	if val, ok := os.LookupEnv("CONFIG"); ok {
		fileName = val
	}

	if fileName != "" {
		var cfg Config
		b, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(b, &cfg)
		if err != nil {
			log.Fatalln(err)
		}

		sup := supportIP{}
		json.Unmarshal(b, &sup)
		if sup.IP != "" {
			_, ip, err := net.ParseCIDR(sup.IP)
			if err == nil {
				c.TrustedIP = ip
			}
		}
	}

	s := flag.String("t", "", "trusted IP")
	if s != nil && *s != "" {
		_, ip, err := net.ParseCIDR(*s)
		if err == nil {
			c.TrustedIP = ip
		}
	}

	flag.StringVar(&c.Addr, "a", "127.0.0.1:8080",
		"Адрес для запуска приложения")

	flag.StringVar(&c.BaseHost, "b", "http://127.0.0.1:8080",
		"Базовый адрес сокращенного URL")

	flag.StringVar(&c.FileStoragePath, "f", "/tmp/short-url-db.json",
		"Путь до файла-хранилища")

	flag.StringVar(&c.DataBaseDSN, "d", "",
		"Параметры БД")

	flag.BoolVar(&c.HTTPS, "s", false, "Использовать HTTPS")

	flag.Parse()

	if val, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
		c.Addr = val
	}

	if val, ok := os.LookupEnv("BASE_URL"); ok {
		c.BaseHost = val
	}

	if val, ok := os.LookupEnv("FILE_STORAGE_PATH"); ok {
		c.FileStoragePath = val
	}

	if val, ok := os.LookupEnv("DATABASE_DSN"); ok {
		c.DataBaseDSN = val
	}

	if val, ok := os.LookupEnv("ENABLE_HTTPS"); ok {
		v, err := strconv.ParseBool(val)
		if err != nil {
			c.HTTPS = false
		}
		c.HTTPS = v
	}

	if val, ok := os.LookupEnv("TRUSTED_SUBNET"); ok {
		_, ip, err := net.ParseCIDR(val)
		if err == nil {
			c.TrustedIP = ip
		}
	}
}

// GetAddress return address
func (c *Config) GetAddress() string {
	return c.Addr
}

// GetShorterURL return prefix of shorten url
func (c *Config) GetShorterURL() string {
	return c.BaseHost + "/"
}

// GetFileStoragePath return file storage path
func (c *Config) GetFileStoragePath() string {
	return c.FileStoragePath
}

// GetDataBaseDSN return database dsn
func (c *Config) GetDataBaseDSN() string {
	return c.DataBaseDSN
}

// GetDataBaseStatus return status of db
func (c *Config) GetDataBaseStatus() bool {
	return c.DataBaseStatus
}

// SetDataBaseStatus set db status
func (c *Config) SetDataBaseStatus(status bool) {
	c.DataBaseStatus = status
}

// GetHTTPS return bool value means should be use HTTPS or HTTP
func (c *Config) GetHTTPS() bool {
	return c.HTTPS
}

// GetTrustedIP return trusted IP
func (c *Config) GetTrustedIP() *net.IPNet {
	return c.TrustedIP
}
