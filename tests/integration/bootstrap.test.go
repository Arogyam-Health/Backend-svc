package integration

import (
	"net/http"
	"os"
	"testing"
	"time"

	"backend-service/internal/bootstrap"
	"backend-service/internal/token"
	"backend-service/tests/dummy"
	"backend-service/tests/helpers"
)

func TestBootstrapLoadsFromPostgres(t *testing.T) {
	srv := dummy.StartDummyServer()
	defer srv.Close()

	db := helpers.SetupTestDB(t)
	helpers.InsertToken(db, "DB_TOKEN", time.Now().Add(30*24*time.Hour))

	rt := token.NewRuntime()

	err := bootstrap.InitToken(
		rt,
		db,
		&http.Client{},
		"test_token.json",
	)

	if err != nil {
		t.Fatal(err)
	}

	if rt.Get() != "DB_TOKEN" {
		t.Fatalf("expected DB_TOKEN, got %s", rt.Get())
	}

	if _, err := os.Stat("test_token.json"); err != nil {
		t.Fatal("token.json not created")
	}
}
