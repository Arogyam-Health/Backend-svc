package helpers

import (
	"database/sql"
	"time"
)

func InsertToken(db *sql.DB, token string, expiry time.Time) error {
	_, err := db.Exec(`
		INSERT INTO instagram_tokens (id, access_token, expires_at)
		VALUES (TRUE, $1, $2)
		ON CONFLICT (id)
		DO UPDATE SET
			access_token = EXCLUDED.access_token,
			expires_at = EXCLUDED.expires_at
	`, token, expiry)

	return err
}
