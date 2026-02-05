package integration

import (
	"net/http"
	"os"
	"testing"
	"time"

	"backend-service/internal/bootstrap"
	"backend-service/internal/instagram"
	"backend-service/internal/token"
	"backend-service/tests/dummy"
	"backend-service/tests/helpers"
)

func TestTokenRefreshBeforeExpiryIntegration(t *testing.T) {
	// Configure environment to use dummy server
	os.Setenv("FB_API_BASE_URL", "http://localhost:9090")
	os.Setenv("APP_ID", "test_app_id")
	os.Setenv("APP_SECRET", "test_app_secret")
	
	srv := dummy.StartDummyServer()
	defer srv.Close()

	redisClient := helpers.SetupTestRedis(t)
	// Insert a token that will expire in 2 days (within 7 days threshold)
	t.Log("Inserting OLD_TOKEN into Redis")
	err := helpers.InsertTokenRedis(redisClient, "OLD_TOKEN", time.Now().Add(2*24*time.Hour))
	if err != nil {
		t.Fatalf("failed to insert token: %v", err)
	}

	rt := token.NewRuntime()
	bootstrap.InitToken(rt, redisClient, &http.Client{}, "test_token.json")

	// Verify bootstrap loaded the OLD_TOKEN
	if rt.Get() != "OLD_TOKEN" {
		t.Fatalf("bootstrap did not load OLD_TOKEN, got %s", rt.Get())
	}

	refreshFn := func() {
		if rt.IsValid() {
			newToken, err := instagram.RefreshAccessToken(&http.Client{}, rt.Get())
			if err != nil {
				t.Fatalf("refresh failed: %v", err)
			}
			rt.Set(newToken)
			token.SaveToRedis(redisClient, newToken)
			t.Logf("Refreshed token: %s", newToken.AccessToken)
		} else {
			t.Log("Token is not valid for refresh")
		}
	}

	refreshFn()

	gotToken := rt.Get()
	if gotToken != "REFRESHED_TOKEN" {
		t.Fatalf("token was not refreshed, expected REFRESHED_TOKEN, got %s", gotToken)
	}
}
