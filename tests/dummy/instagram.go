package dummy

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func StartDummyServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/refresh_access_token", func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("access_token")
		if len(accessToken) < len("NEW_DUMMY_TOKEN") || accessToken[:len("NEW_DUMMY_TOKEN")] != "NEW_DUMMY_TOKEN" {
			http.Error(w, `{"error": "invalid access_token"}`, http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "NEW_DUMMY_TOKEN_V" + fmt.Sprintf("%d", 1000+rand.Intn(9000)),
			"expires_in":   5184000,
		})
	})

	mux.HandleFunc("/17841400000000000/media", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": []map[string]string{
				{
					"id":         "1",
					"media_type": "IMAGE",
					"media_url":  "https://example.com/image.jpg",
					"permalink":  "https://instagram.com/p/abc",
				},
			},
		})
	})

	srv := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	return srv
}
