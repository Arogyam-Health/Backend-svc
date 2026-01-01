package config

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func ConnectPostgres() (*sql.DB, error) {
	dbURL := os.Getenv("DATABASE_URL")
	return sql.Open("postgres", dbURL)
}
