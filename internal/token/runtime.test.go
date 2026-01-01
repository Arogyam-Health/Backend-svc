package token

import (
	"testing"
	"time"
)

func TestRuntimeSetAndGet(t *testing.T) {
	rt := NewRuntime()

	token := Token{
		AccessToken: "TEST_TOKEN",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}

	rt.Set(token)

	if rt.Get() != "TEST_TOKEN" {
		t.Fatalf("expected TEST_TOKEN, got %s", rt.Get())
	}
}

func TestTokenExpiresSoon(t *testing.T) {
	rt := NewRuntime()

	rt.Set(Token{
		AccessToken: "EXPIRING",
		ExpiresAt:   time.Now().Add(3 * 24 * time.Hour),
	})

	if !rt.IsValid() {
		t.Fatal("expected token to expire soon")
	}
}

func TestTokenNotExpiringSoon(t *testing.T) {
	rt := NewRuntime()

	rt.Set(Token{
		AccessToken: "SAFE",
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour),
	})

	if rt.IsValid() {
		t.Fatal("token should not be expiring soon")
	}
}
