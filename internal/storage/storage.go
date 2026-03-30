package storage

import (
	"context"
	"os"
	"strings"
)

type APIKeyStorage interface {
	Get(ctx context.Context, name string) (token string, ok bool)
}

type APIKeyENVStorage struct{}

func (APIKeyENVStorage) Get(ctx context.Context, name string) (string, bool) {
	_ = ctx
	key := "BOT_API_KEY_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
	v := os.Getenv(key)
	if v == "" {
		return "", false
	}
	return v, true
}
