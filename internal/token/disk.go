package token

import (
	"encoding/json"
	"os"
)

func NewStore() *TokenRuntime {
	return &TokenRuntime{}
}

/* Disk persistence methods for TokenRuntime */

func LoadFromDisk(path string) (*Token, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var token Token
	err = json.Unmarshal(data, &token)
	return &token, err
}

func SaveToDisk(path string, token *Token) error {
	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644) // rw-r--r--
}
