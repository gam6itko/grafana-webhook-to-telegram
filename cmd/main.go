package main

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/gam6itko/grafana-webhook-to-telegram/internal/config"
	"github.com/gam6itko/grafana-webhook-to-telegram/internal/handler"
	"github.com/gam6itko/grafana-webhook-to-telegram/internal/storage"
	"github.com/gam6itko/grafana-webhook-to-telegram/internal/telegram"
)

func main() {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		panic(err)
	}
	logger := initLogger(&cfg.Logs)
	defer func() { _ = logger.Sync() }()
	tgClient := telegram.NewClient(telegram.WithBaseURL(cfg.TelegramAPIHost))
	h := handler.NewWebhook(logger, storage.APIKeyENVStorage{}, tgClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/{bot_name}/{chat_id}", h.ServeHTTP)
	mux.HandleFunc("PUT /api/{bot_name}/{chat_id}", h.ServeHTTP)

	logger.Info("server starting", zap.String("addr", cfg.ListenAddr))
	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
}

// initLogger builds a zap.Logger from logs config (mode, level, disable_stacktrace).
func initLogger(logs *config.LogsConfig) *zap.Logger {
	var cfg zap.Config
	if logs.Mode == "development" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}
	if logs.Level != "" {
		var l zap.AtomicLevel
		if err := l.UnmarshalText([]byte(logs.Level)); err == nil {
			cfg.Level = l
		}
	}
	cfg.DisableStacktrace = logs.DisableStacktrace
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	return logger
}
