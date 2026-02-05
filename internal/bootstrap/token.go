package bootstrap

import (
	"log"
	"net/http"
	"time"

	"backend-service/internal/instagram"
	"backend-service/internal/token"

	"github.com/redis/go-redis/v9"
)

func InitToken(
	runtime *token.TokenRuntime,
	redisClient *redis.Client,
	client *http.Client,
	tokenPath string,
) error {

	// 1. Try disk
	if t, err := token.LoadFromDisk(tokenPath); err == nil {
		if time.Now().Before(t.ExpiresAt) {
			runtime.Set(*t)
			return nil
		} else {
			log.Printf("[BOOTSTRAP] Token from disk is expired - %v", t.ExpiresAt)
		}
	} else {
		log.Printf("[BOOTSTRAP] No token found on disk: %v", err)
	}

	// 2. Try Redis
	if t, err := token.LoadFromRedis(redisClient); err == nil {
		if time.Now().Before(t.ExpiresAt) {
			runtime.Set(*t)
			token.SaveToDisk(tokenPath, t)
			return nil
		} else {
			log.Printf("[BOOTSTRAP] Token from Redis is expired - %v", t.ExpiresAt)
		}
	} else {
		log.Printf("[BOOTSTRAP] No token found in Redis: %v", err)
	}

	// 3. Refresh from Instagram
	log.Println("[BOOTSTRAP] Refreshing token from Instagram...")
	newInstagramToken, err := instagram.RefreshAccessToken(client, runtime.Get())
	if err != nil {
		return err
	}

	// Map instagram.Token to token.Token
	newToken := token.Token{
		AccessToken: newInstagramToken.AccessToken,
		ExpiresAt:   newInstagramToken.ExpiresAt,
		// Add other fields as necessary
	}

	log.Printf("[BOOTSTRAP] New token expires at: %v", newToken.ExpiresAt)
	log.Println("[BOOTSTRAP] Storing new token...")
	runtime.Set(newToken)
	token.SaveToDisk(tokenPath, &newToken)
	token.SaveToRedis(redisClient, newToken)

	return nil
}
