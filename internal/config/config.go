package config

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

const (
	DefaultListenAddr   = "127.0.0.1:8080"
	DefaultTelegramHost = "https://api.telegram.org"
)

type LogsConfig struct {
	Mode              string `env:"MODE" envDefault:"development"`
	Level             string `env:"LEVEL" envDefault:"info"`
	DisableStacktrace bool   `env:"DISABLE_STACKTRACE" envDefault:"true"`
}

type Config struct {
	ListenAddr      string     `env:"HTTP_SERVER_LISTEN_ADDR" envDefault:"127.0.0.1:8080"`
	TelegramAPIHost string     `env:"TELEGRAM_API_HOST" envDefault:"https://api.telegram.org"`
	Logs            LogsConfig `envPrefix:"LOG_"`
}

func LoadFromEnv() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	cfg.TelegramAPIHost = strings.TrimRight(cfg.TelegramAPIHost, "/")
	return cfg, nil
}
