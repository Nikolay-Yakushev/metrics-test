package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)


type Config struct {
	LogLevel   string `env:"LOG_LEVEL" envDefault:"INFO"`
	ServerPort int    `env:"SERVER_PORT" envDefault:"8080"`
	ServerHost string `env:"SERVER_PORT" envDefault:"127.0.0.1"`
}


func New() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &cfg
}