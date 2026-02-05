package helpers

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"backend-service/internal/token"

	"github.com/redis/go-redis/v9"
)

func InsertTokenRedis(client *redis.Client, accessToken string, expiry time.Time) error {
	tok := token.Token{
		AccessToken: accessToken,
		ExpiresAt:   expiry,
	}

	data, err := json.Marshal(tok)
	if err != nil {
		return err
	}

	ttl := time.Until(expiry)
	if ttl < 0 {
		ttl = 0
	}

	// Use test-specific key from environment to avoid affecting production
	testTokenKey := os.Getenv("REDIS_TOKEN_KEY")
	if testTokenKey == "" {
		testTokenKey = "test_instagram_token" // fallback
	}

	return client.Set(context.Background(), testTokenKey, data, ttl).Err()
}
