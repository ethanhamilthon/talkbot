package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	telegramToken string
	redisAddr string
}

func New() *Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	return &Config{
		telegramToken: os.Getenv("TGAPI"),
		redisAddr:  os.Getenv("REDIS_ADDR"),
	}
}

func (c *Config) GetTelegramToken() string {
	return c.telegramToken
}

func (c *Config) GetRedisAddr() string {
	return c.redisAddr
}

