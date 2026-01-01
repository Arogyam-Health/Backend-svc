package bootstrap

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"backend-service/internal/instagram"
	"backend-service/internal/token"
)

func InitToken(
	runtime *token.TokenRuntime,
	db *sql.DB,
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

	// 2. Try Postgres
	if t, err := token.LoadFromDB(db); err == nil {
		if time.Now().Before(t.ExpiresAt) {
			runtime.Set(*t)
			token.SaveToDisk(tokenPath, t)
			return nil
		} else {
			log.Printf("[BOOTSTRAP] Token from database is expired - %v", t.ExpiresAt)
		}
	} else {
		log.Printf("[BOOTSTRAP] No token found in database: %v", err)
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
	token.SaveToDB(db, newToken)

	return nil
}
