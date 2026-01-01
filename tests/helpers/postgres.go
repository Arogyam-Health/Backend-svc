package helpers

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}

	// ensure table exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS instagram_tokens (
			id BOOLEAN PRIMARY KEY DEFAULT TRUE,
			access_token TEXT NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ DEFAULT now()
		)
	`)
	if err != nil {
		t.Fatal(err)
	}

	// clean state before each test
	_, err = db.Exec(`DELETE FROM instagram_tokens`)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}
