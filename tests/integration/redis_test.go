package integration

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"backend-service/internal/token"
	"backend-service/tests/helpers"
)

var ctx = context.Background()

func getTestTokenKey() string {
	if key := os.Getenv("REDIS_TOKEN_KEY"); key != "" {
		return key
	}
	return "test_instagram_token" // fallback
}

func TestUpstashRedisReadWrite(t *testing.T) {
	redisClient := helpers.SetupTestRedis(t)
	testTokenKey := getTestTokenKey()

	// Test 1: Write token to Redis
	t.Run("WriteToken", func(t *testing.T) {
		testToken := token.Token{
			AccessToken: "test_access_token_12345",
			ExpiresAt:   time.Now().Add(60 * 24 * time.Hour), // 60 days from now
		}

		err := token.SaveToRedis(redisClient, testToken)
		if err != nil {
			t.Fatalf("failed to save token to Redis: %v", err)
		}

		// Verify token was saved by checking Redis directly
		val, err := redisClient.Get(ctx, testTokenKey).Result()
		if err != nil {
			t.Fatalf("failed to retrieve token from Redis: %v", err)
		}

		var storedToken token.Token
		if err := json.Unmarshal([]byte(val), &storedToken); err != nil {
			t.Fatalf("failed to unmarshal token: %v", err)
		}

		if storedToken.AccessToken != testToken.AccessToken {
			t.Errorf("expected token %s, got %s", testToken.AccessToken, storedToken.AccessToken)
		}

		// Check if TTL is set correctly (should be approximately 60 days)
		ttl := redisClient.TTL(ctx, testTokenKey).Val()
		expectedTTL := 60 * 24 * time.Hour
		if ttl < expectedTTL-time.Hour || ttl > expectedTTL+time.Hour {
			t.Errorf("TTL not set correctly: expected ~%v, got %v", expectedTTL, ttl)
		}
	})

	// Test 2: Read token from Redis
	t.Run("ReadToken", func(t *testing.T) {
		// Insert a token directly using helper
		expectedToken := "read_test_token_67890"
		expiry := time.Now().Add(30 * 24 * time.Hour)

		err := helpers.InsertTokenRedis(redisClient, expectedToken, expiry)
		if err != nil {
			t.Fatalf("failed to insert token: %v", err)
		}

		// Read the token using LoadFromRedis
		loadedToken, err := token.LoadFromRedis(redisClient)
		if err != nil {
			t.Fatalf("failed to load token from Redis: %v", err)
		}

		if loadedToken.AccessToken != expectedToken {
			t.Errorf("expected token %s, got %s", expectedToken, loadedToken.AccessToken)
		}

		// Verify the expiry time is close to what we set (within 1 second tolerance)
		timeDiff := loadedToken.ExpiresAt.Sub(expiry)
		if timeDiff < -time.Second || timeDiff > time.Second {
			t.Errorf("expiry time mismatch: expected %v, got %v", expiry, loadedToken.ExpiresAt)
		}
	})

	// Test 3: Update existing token
	t.Run("UpdateToken", func(t *testing.T) {
		// Insert initial token
		initialToken := token.Token{
			AccessToken: "initial_token",
			ExpiresAt:   time.Now().Add(10 * 24 * time.Hour),
		}
		token.SaveToRedis(redisClient, initialToken)

		// Update with new token
		updatedToken := token.Token{
			AccessToken: "updated_token",
			ExpiresAt:   time.Now().Add(50 * 24 * time.Hour),
		}
		err := token.SaveToRedis(redisClient, updatedToken)
		if err != nil {
			t.Fatalf("failed to update token: %v", err)
		}

		// Verify the update
		loadedToken, err := token.LoadFromRedis(redisClient)
		if err != nil {
			t.Fatalf("failed to load updated token: %v", err)
		}

		if loadedToken.AccessToken != updatedToken.AccessToken {
			t.Errorf("expected updated token %s, got %s", updatedToken.AccessToken, loadedToken.AccessToken)
		}
	})

	// Test 4: Handle expired token (TTL = 0)
	t.Run("ExpiredToken", func(t *testing.T) {
		expiredToken := token.Token{
			AccessToken: "expired_token",
			ExpiresAt:   time.Now().Add(-1 * time.Hour), // Already expired
		}

		err := token.SaveToRedis(redisClient, expiredToken)
		if err != nil {
			t.Fatalf("failed to save expired token: %v", err)
		}

		// The token should still be saved but with TTL=0 (expires immediately)
		// Try to load it - it might not exist anymore
		_, err = token.LoadFromRedis(redisClient)
		// We expect either a successful load or redis.Nil error
		if err != nil {
			t.Logf("expired token could not be loaded (expected): %v", err)
		}
	})

	// Test 5: Read non-existent token
	t.Run("NonExistentToken", func(t *testing.T) {
		// Clear any existing token
		redisClient.Del(ctx, testTokenKey)

		_, err := token.LoadFromRedis(redisClient)
		if err == nil {
			t.Fatal("expected error when reading non-existent token, got nil")
		}
	})
}
