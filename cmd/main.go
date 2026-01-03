package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"backend-service/api"
	"backend-service/internal/bootstrap"
	"backend-service/internal/cache"
	"backend-service/internal/config"
	"backend-service/internal/instagram"
	"backend-service/internal/scheduler"
	"backend-service/internal/token"
	"backend-service/middleware"
)

func main() {
	cfg := config.LoadConfig()
	store := cache.NewStore()

	db, err := config.ConnectPostgres()
	if err != nil {
		log.Fatal("failed to connect postgres:", err)
	}
	defer db.Close()

	runtimeToken := token.NewRuntime()

	client := instagram.NewClient()

	// Bootstrap
	if err := bootstrap.InitToken(
		runtimeToken,
		db,
		client,
		"token.json",
	); err != nil {
		log.Fatal("[BOOTSTRAP] ", err)
	}

	service := instagram.Service{
		Client:     client,
		IgUserID:   cfg.IgUserID,
		TokenStore: runtimeToken,
	}

	syncFn := func() {
		const maxAttempts = 3
		for i := 1; i <= maxAttempts; i++ {
			media, err := service.FetchMedia()
			if err == nil {
				store.SetMedia(media)
				return
			}
			log.Printf("[MEDIA] Attempt %d/%d failed: %v", i, maxAttempts, err)
			if i < maxAttempts {
				time.Sleep(time.Duration(i) * time.Second) // simple backoff: 1s, 2s, ...
			}
		}
		log.Printf("[MEDIA] All attempts to fetch media failed")
	}

	refTok := func() {
		if runtimeToken.IsValid() {
			newToken, err := instagram.RefreshAccessToken(client, runtimeToken.Get())
			if err != nil {
				log.Printf("failed to refresh token: %v", err)
				return
			}
			runtimeToken.Set(newToken)
			token.SaveToDisk("token.json", &newToken)
			token.SaveToDB(db, newToken)

			log.Printf("[TOKEN] refreshed access token")
		} else {
			log.Printf("[TOKEN] access token is still valid, no need to refresh")
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	scheduler.Start(ctx, syncFn, refTok)
	defer cancel()

	mux := http.NewServeMux()
	
	mux.HandleFunc("/media", api.MediaHandler(store))
	mux.HandleFunc("/media/getIdsOnly", api.MediaIdsHandler(store))
	mux.HandleFunc("/ready", api.ReadyHandler)

	handler := middleware.CORS(mux)

	log.Printf("Started server on %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
