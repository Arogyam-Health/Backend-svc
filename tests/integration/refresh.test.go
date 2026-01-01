package integration

import (
	"net/http"
	"testing"
	"time"

	"backend-service/internal/bootstrap"
	"backend-service/internal/instagram"
	"backend-service/internal/token"
	"backend-service/tests/dummy"
	"backend-service/tests/helpers"
)

func TestTokenRefreshBeforeExpiry(t *testing.T) {
	srv := dummy.StartDummyServer()
	defer srv.Close()

	db := helpers.SetupTestDB(t)
	helpers.InsertToken(db, "OLD_TOKEN", time.Now().Add(2*24*time.Hour))

	rt := token.NewRuntime()
	bootstrap.InitToken(rt, db, &http.Client{}, "test_token.json")

	refreshFn := func() {
		if rt.IsValid() {
			newToken, _ := instagram.RefreshAccessToken(&http.Client{}, rt.Get())
			rt.Set(newToken)
			token.SaveToDB(db, newToken)
		}
	}

	refreshFn()

	if rt.Get() != "REFRESHED_TOKEN" {
		t.Fatal("token was not refreshed")
	}
}
