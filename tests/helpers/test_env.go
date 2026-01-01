package helpers

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".env.test")
	m.Run()
}
