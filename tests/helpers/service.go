package helpers

import (
	"net/http"

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

	return &instagram.Service{
		Client:     client,
		IgUserID:   igUserID,
		TokenStore: rt,
	}
}
