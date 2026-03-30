package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gam6itko/grafana-webhook-to-telegram/internal/handler"
)

var _ handler.MessageSender = (*Client)(nil)

type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient returns a Client with optional configuration. Defaults: http.DefaultClient.
func NewClient(opts ...Option) *Client {
	c := &Client{
		http:    http.DefaultClient,
		baseURL: "https://api.telegram.org",
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type sendMessageBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type apiResponse struct {
	OK          bool   `json:"ok"`
	Description string `json:"description"`
}

func (c *Client) SendMessage(ctx context.Context, token, chatID, text string) error {
	if c.http == nil {
		c.http = http.DefaultClient
	}
	base := strings.TrimRight(c.baseURL, "/")
	url := fmt.Sprintf("%s/bot%s/sendMessage", base, token)

	payload, err := json.Marshal(sendMessageBody{ChatID: chatID, Text: text})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram http %d: %s", resp.StatusCode, string(body))
	}

	var ar apiResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return fmt.Errorf("telegram response: %w", err)
	}
	if !ar.OK {
		return fmt.Errorf("telegram: %s", ar.Description)
	}
	return nil
}
