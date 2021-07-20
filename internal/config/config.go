package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ConnectionString string `env:"DB_URI,required"`
	Port             int    `env:"PORT" envDefault:"8080"`
	Chats            []int  `env:"CHATS" envSeparator:","`
	AppID            int    `env:"APP_ID,required"`
	AppHash          string `env:"APP_HASH,required"`
	Phone            string `env:"PHONE,required"`
	Password         string `env:"PASSWORD"`
	N                int    `env:"N"` // number of messages to save
}

func Parse() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
