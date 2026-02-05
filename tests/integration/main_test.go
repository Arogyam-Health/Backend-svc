package integration

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	// Load .env.test for integration tests
	_ = godotenv.Load("../../.env.test")

	// REDIS_TOKEN_KEY is now set from .env.test, no need to override
	// unless .env.test doesn't have it
	if os.Getenv("REDIS_TOKEN_KEY") == "" {
		os.Setenv("REDIS_TOKEN_KEY", "test_instagram_token")
	}

	os.Exit(m.Run())
}
