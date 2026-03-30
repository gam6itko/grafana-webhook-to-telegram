package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/gam6itko/grafana-webhook-to-telegram/internal/storage"
	"go.uber.org/zap"
)

var botNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type MessageSender interface {
	SendMessage(ctx context.Context, token, chatID, text string) error
}

type GrafanaWebhook struct {
	Message string `json:"message"`
	Title   string `json:"title"`
	Status  string `json:"status"`
}

type Webhook struct {
	log    *zap.Logger
	keys   storage.APIKeyStorage
	sender MessageSender
}

func NewWebhook(log *zap.Logger, keys storage.APIKeyStorage, sender MessageSender) *Webhook {
	return &Webhook{log: log, keys: keys, sender: sender}
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	botName := r.PathValue("bot_name")
	chatID := r.PathValue("chat_id")

	if !botNamePattern.MatchString(botName) {
		http.NotFound(w, r)
		return
	}

	var payload GrafanaWebhook
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.log.Warn("invalid json body", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: "invalid json"})
		return
	}

	text := payload.Message
	if text == "" {
		text = payload.Title
	}

	h.log.Info("incoming webhook",
		zap.String("bot_name", botName),
		zap.String("grafana.title", payload.Title),
		zap.String("grafana.status", payload.Status),
	)

	token, ok := h.keys.Get(r.Context(), botName)
	if !ok {
		h.log.Error("bot not found", zap.String("bot_name", botName))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Error: "bot with name " + botName + " not found",
		})
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	if err := h.sender.SendMessage(ctx, token, chatID, text); err != nil {
		h.log.Error("telegram send failed", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
