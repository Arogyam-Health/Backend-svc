package token

import (
	"os"
	"testing"
	"time"
)

func TestSaveAndLoadTokenFromDisk(t *testing.T) {
	path := "test_token.json"
	defer os.Remove(path)

	expected := Token{
		AccessToken: "DISK_TOKEN",
		ExpiresAt:   time.Now().Add(10 * time.Hour),
	}

	err := SaveToDisk(path, &expected)
	if err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadFromDisk(path)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.AccessToken != expected.AccessToken {
		t.Fatalf("expected %s, got %s", expected.AccessToken, loaded.AccessToken)
	}
}
