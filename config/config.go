package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type (
	Config struct {
		Postgres Postgres `json:"postgres"`
		Redis    Redis    `json:"redis"`
		Telegram Telegram `json:"telegram"`
	}

	Postgres struct {
		URL string `json:"url"`
	}

	Redis struct {
		Password     string `json:"password"`
		Host         string `json:"host"`
		Db           int    `json:"db"`
		MinIdleConns int    `json:"min_idle_conns"`
	}

	Telegram struct {
		Token string `json:"token"`
	}
)

func New() (*Config, error) {
	_, err := os.Stat("env/bot.env")
	if err == nil {
		err = godotenv.Load("env/bot.env")
		if err != nil {
			return nil, err
		}
	}

	config := &Config{
		Postgres: Postgres{
			URL: os.Getenv("POSTGRES_URL"),
		},
		Redis: Redis{
			Password:     os.Getenv("REDIS_PASSWORD"),
			Host:         os.Getenv("REDIS_HOST"),
			Db:           parseEnvInt(os.Getenv("REDIS_DB")),
			MinIdleConns: parseEnvInt(os.Getenv("REDIS_MIN_IDLE_CONNS")),
		},
		Telegram: Telegram{
			Token: os.Getenv("TG_TOKEN"),
		},
	}

	return config, nil
}

func parseEnvInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}
