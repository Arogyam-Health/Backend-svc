package config

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ConnectRedis() (*redis.Client, error) {
	// Get Redis URL from environment
	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		return nil, fmt.Errorf("UPSTASH_REDIS_REST_URL or UPSTASH_REDIS_REST_TOKEN not set")
	}

	// Convert https:// to rediss:// for Upstash
	if len(redisURL) > 8 && redisURL[:8] == "https://" {
		redisURL = "rediss://" + redisURL[8:]
	}

	// Parse the URL
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Set the password from token
	opt.Password = redisToken

	client := redis.NewClient(opt)

	// Test the connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
