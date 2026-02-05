package dummy

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

func StartDummyServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/oauth/access_token", func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("fb_exchange_token")
		// For testing: return REFRESHED_TOKEN for OLD_TOKEN
		if accessToken == "OLD_TOKEN" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "REFRESHED_TOKEN",
				"expires_in":   5184000,
			})
			return
		}
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

	mux.HandleFunc("/refresh_access_token", func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.URL.Query().Get("access_token")
		// For testing: return REFRESHED_TOKEN for OLD_TOKEN
		if accessToken == "OLD_TOKEN" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "REFRESHED_TOKEN",
				"expires_in":   5184000,
			})
			return
		}
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

	// Generic handler for any user ID
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle any media endpoint
		if r.URL.Path != "/" && r.URL.Path != "" {
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
			return
		}
		http.NotFound(w, r)
	})

	srv := &http.Server{
		Addr:    ":9090",
		Handler: mux,
	}

	go srv.ListenAndServe()

	return srv
}
