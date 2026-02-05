package token

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func getTokenKey() string {
	if key := os.Getenv("REDIS_TOKEN_KEY"); key != "" {
		return key
	}
	return "instagram_token"
}

func LoadFromRedis(client *redis.Client) (*Token, error) {
	val, err := client.Get(ctx, getTokenKey()).Result()
	if err != nil {
		return nil, err
	}

	var t Token
	if err := json.Unmarshal([]byte(val), &t); err != nil {
		return nil, err
	}

	return &t, nil
}

func SaveToRedis(client *redis.Client, t Token) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	// Calculate TTL based on token expiration
	ttl := time.Until(t.ExpiresAt)
	if ttl < 0 {
		ttl = 0
	}

	return client.Set(ctx, getTokenKey(), data, ttl).Err()
}
