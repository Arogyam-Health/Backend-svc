package integration

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"backend-service/internal/bootstrap"
	"backend-service/internal/token"
	"backend-service/tests/dummy"
	"backend-service/tests/helpers"
)

func TestBootstrapInitTokenFromRedis(t *testing.T) {
	srv := dummy.StartDummyServer()
	defer srv.Close()

	redisClient := helpers.SetupTestRedis(t)
	// Use timestamp to create unique token for this test
	testToken := fmt.Sprintf("REDIS_TOKEN_%d", time.Now().UnixNano())
	helpers.InsertTokenRedis(redisClient, testToken, time.Now().Add(30*24*time.Hour))

	rt := token.NewRuntime()

	err := bootstrap.InitToken(
		rt,
		redisClient,
		&http.Client{},
		"test_token.json",
	)

	if err != nil {
		t.Fatal(err)
	}

	if rt.Get() != testToken {
		t.Fatalf("expected %s, got %s", testToken, rt.Get())
	}

	if _, err := os.Stat("test_token.json"); err != nil {
		t.Fatal("token.json not created")
	}
}
