package helpers

import (
	"net/http"
	"os"

	"backend-service/internal/instagram"
	"backend-service/internal/token"
)

func NewTestService(igUserID string) *instagram.Service {
	rt := token.NewRuntime()

	// default test token
	rt.Set(token.Token{
		AccessToken: "TEST_TOKEN",
	})

	client := &http.Client{}

	// Point to dummy server for tests
	os.Setenv("FB_API_BASE_URL", "http://localhost:9090")

	return &instagram.Service{
		Client:     client,
		IgUserID:   igUserID,
		TokenStore: rt,
	}
}
