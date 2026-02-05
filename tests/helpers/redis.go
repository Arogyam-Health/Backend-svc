package helpers

import (
	"context"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func SetupTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		t.Fatal("UPSTASH_REDIS_REST_URL or UPSTASH_REDIS_REST_TOKEN not set")
	}

	// Convert https:// to rediss:// for Upstash
	if len(redisURL) > 8 && redisURL[:8] == "https://" {
		redisURL = "rediss://" + redisURL[8:]
	}

	// Parse the URL with token
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Fatal(err)
	}

	// Set the password from token
	opt.Password = redisToken

	client := redis.NewClient(opt)

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("failed to connect to Redis: %v", err)
	}

	// Use test-specific key from environment to avoid affecting production
	testTokenKey := os.Getenv("REDIS_TOKEN_KEY")
	if testTokenKey == "" {
		testTokenKey = "test_instagram_token" // fallback
	}

	// Clean state before each test - ensure complete deletion
	delCount, err := client.Del(ctx, testTokenKey).Result()
	if err != nil {
		t.Logf("Warning: failed to delete %s: %v", testTokenKey, err)
	} else {
		t.Logf("Deleted %d keys from Redis", delCount)
	}

	// Also clean up test token file
	os.Remove("test_token.json")

	t.Cleanup(func() {
		client.Del(ctx, testTokenKey)
		os.Remove("test_token.json")
		client.Close()
	})

	return client
}
