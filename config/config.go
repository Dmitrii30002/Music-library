package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST     string `env:"DB_HOST" required:"true"`
	DB_PORT     string `env:"DB_PORT" required:"true"`
	DB_USER     string `env:"DB_USER" required:"true"`
	DB_PASSWORD string `env:"DB_PASSWORD" required:"true"`
	DB_NAME     string `env:"DB_NAME" required:"true"`
	Port        string `env:"Port" required:"true"`
	DB_DSN      string `env:"DB_DSN" required:"true"`
}

var (
	cfg  Config
	once sync.Once
)

func Get() Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("[ERROR] Failed to load config: %v", err)
		}
		cfg = Config{
			DB_HOST:     os.Getenv("DB_DSN"),
			DB_PORT:     os.Getenv("DB_PORT"),
			DB_USER:     os.Getenv("DB_USER"),
			DB_PASSWORD: os.Getenv("DB_PASSWORD"),
			DB_NAME:     os.Getenv("DB_NAME"),
			Port:        os.Getenv("Port"),
			DB_DSN:      os.Getenv("DB_DSN"),
		}
	})
	return cfg
}
