package storage

import (
	"context"
	"testing"
)

func TestAPIKeyENVStorage_Get(t *testing.T) {
	ctx := context.Background()
	s := APIKeyENVStorage{}

	t.Run("lowercase name maps to upper env key", func(t *testing.T) {
		t.Setenv("BOT_API_KEY_FOO", "token-abc")
		got, ok := s.Get(ctx, "foo")
		if !ok || got != "token-abc" {
			t.Fatalf("Get(foo) = %q, %v; want token-abc, true", got, ok)
		}
	})

	t.Run("hyphen in name becomes underscore in env key", func(t *testing.T) {
		t.Setenv("BOT_API_KEY_MY_BOT", "tok")
		got, ok := s.Get(ctx, "my-bot")
		if !ok || got != "tok" {
			t.Fatalf("Get(my-bot) = %q, %v; want tok, true", got, ok)
		}
	})

	t.Run("missing env returns false", func(t *testing.T) {
		t.Setenv("BOT_API_KEY_MISSINGX", "")
		got, ok := s.Get(ctx, "missingx")
		if ok || got != "" {
			t.Fatalf("Get(missingx) = %q, %v; want empty, false", got, ok)
		}
	})
}
