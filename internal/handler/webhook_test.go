package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gam6itko/grafana-webhook-to-telegram/internal/storage"
	"go.uber.org/zap"
)

type mockSender struct {
	err error
}

func (m *mockSender) SendMessage(ctx context.Context, token, chatID, text string) error {
	return m.err
}

func TestMaskTokenInPath(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"/bot123456:ABCDefgh/sendMessage", "/bot***/sendMessage"},
		{"/bot123456:ABCDefgh/getMe", "/bot***/getMe"},
		{"/other/path", "/other/path"},
	}
	for _, c := range cases {
		if got := maskTokenInPath(c.in); got != c.want {
			t.Errorf("maskTokenInPath(%q) = %q; want %q", c.in, got, c.want)
		}
	}
}

func TestWebhook_ServeHTTP(t *testing.T) {
	log := zap.NewNop()
	body := `{"message":"hello","title":"t","status":"firing"}`

	t.Run("invalid bot name returns 404", func(t *testing.T) {
		mux := http.NewServeMux()
		h := NewWebhook(log, storage.APIKeyENVStorage{}, &mockSender{})
		mux.HandleFunc("POST /api/{bot_name}/{chat_id}", h.ServeHTTP)

		req := httptest.NewRequest(http.MethodPost, "/api/foo.bar/1", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d; want 404", rec.Code)
		}
	})

	t.Run("unknown bot returns 404 json", func(t *testing.T) {
		t.Setenv("BOT_API_KEY_ALERTS", "")
		mux := http.NewServeMux()
		h := NewWebhook(log, storage.APIKeyENVStorage{}, &mockSender{})
		mux.HandleFunc("POST /api/{bot_name}/{chat_id}", h.ServeHTTP)

		req := httptest.NewRequest(http.MethodPost, "/api/alerts/99", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("status = %d; want 404", rec.Code)
		}
		var er errorResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &er); err != nil {
			t.Fatal(err)
		}
		if er.Error != "bot with name alerts not found" {
			t.Fatalf("error = %q", er.Error)
		}
	})

	t.Run("success returns 204", func(t *testing.T) {
		t.Setenv("BOT_API_KEY_ALERTS", "dummy-token")
		mux := http.NewServeMux()
		h := NewWebhook(log, storage.APIKeyENVStorage{}, &mockSender{})
		mux.HandleFunc("POST /api/{bot_name}/{chat_id}", h.ServeHTTP)

		req := httptest.NewRequest(http.MethodPost, "/api/alerts/-1001", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status = %d; want 204", rec.Code)
		}
	})
}
