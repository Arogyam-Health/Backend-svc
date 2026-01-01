package token

import (
	"database/sql"
)

func LoadFromDB(db *sql.DB) (*Token, error) {
	row := db.QueryRow(`
		SELECT access_token, expires_at
		FROM instagram_tokens
		WHERE id = TRUE
	`)

	var t Token
	err := row.Scan(&t.AccessToken, &t.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func SaveToDB(db *sql.DB, t Token) error {
	_, err := db.Exec(`
		INSERT INTO instagram_tokens (id, access_token, expires_at)
		VALUES (TRUE, $1, $2)
		ON CONFLICT (id)
		DO UPDATE SET
			access_token = EXCLUDED.access_token,
			expires_at = EXCLUDED.expires_at,
			updated_at = now()
	`, t.AccessToken, t.ExpiresAt)

	return err
}
