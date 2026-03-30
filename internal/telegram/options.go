package telegram

import "net/http"

// Option configures Client via NewClient.
type Option func(*Client)

// WithBaseURL sets the Telegram Bot API base URL (e.g. https://api.telegram.org).
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithHTTP sets the HTTP client used for requests. Nil is ignored.
func WithHTTP(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.http = h
		}
	}
}
