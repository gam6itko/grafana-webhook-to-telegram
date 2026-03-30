package config

import (
	"os"
	"strings"
)

const (
	DefaultListenAddr   = "127.0.0.1:8080"
	DefaultTelegramHost = "https://api.telegram.org"
)

type Config struct {
	ListenAddr      string
	TelegramAPIHost string
}

func LoadFromEnv() Config {
	listen := os.Getenv("HTTP_SERVER_LISTEN_ADDR")
	if listen == "" {
		listen = DefaultListenAddr
	}
	tg := os.Getenv("TELEGRAM_API_HOST")
	if tg == "" {
		tg = DefaultTelegramHost
	}
	return Config{
		ListenAddr:      listen,
		TelegramAPIHost: strings.TrimRight(tg, "/"),
	}
}
