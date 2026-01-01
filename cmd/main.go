package main

import (
	"context"
	"log"
	"net/http"

	"backend-service/api"
	"backend-service/internal/bootstrap"
	"backend-service/internal/cache"
	"backend-service/internal/config"
	"backend-service/internal/instagram"
	"backend-service/internal/scheduler"
	"backend-service/internal/token"
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
		media, err := service.FetchMedia()
		if err != nil {
			log.Printf("[MEDIA] Error fetching media: %v", err)
			return
		}
		store.SetMedia(media)
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

	http.HandleFunc("/media", api.MediaHandler(store))
	http.HandleFunc("/ready", api.ReadyHandler)

	log.Printf("Started server on localhost:%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
