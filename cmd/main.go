package main

import (
	"net/http"

	"github.com/gam6itko/grafana-webhook-to-telegram/internal/config"
	"github.com/gam6itko/grafana-webhook-to-telegram/internal/handler"
	"github.com/gam6itko/grafana-webhook-to-telegram/internal/storage"
	"github.com/gam6itko/grafana-webhook-to-telegram/internal/telegram"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg := config.LoadFromEnv()
	tgClient := &telegram.Client{
		BaseURL: cfg.TelegramAPIHost,
		HTTP:    http.DefaultClient,
	}
	h := handler.NewWebhook(logger, storage.APIKeyENVStorage{}, tgClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/{bot_name}/{chat_id}", h.ServeHTTP)
	mux.HandleFunc("PUT /api/{bot_name}/{chat_id}", h.ServeHTTP)

	logger.Info("server starting", zap.String("addr", cfg.ListenAddr))
	if err := http.ListenAndServe(cfg.ListenAddr, mux); err != nil {
		logger.Fatal("listen", zap.Error(err))
	}
}
