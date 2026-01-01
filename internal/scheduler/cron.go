package scheduler

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
)

func mustEnvInt(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid %s: %v", key, err)
	}
	return i
}

func Start(ctx context.Context, syncFn func(), refTok func()) {
	go func() {
		dataTicker := time.NewTicker(time.Duration(mustEnvInt("MEDIA_SYNC_TIME", 45)) * time.Minute)
		tokenTicker := time.NewTicker(time.Duration(mustEnvInt("TOKEN_REFRESH_TIME", 30)) * time.Hour * 24)
		defer dataTicker.Stop()
		defer tokenTicker.Stop()

		// run once on startup
		go syncFn()
		go refTok()

		for {
			select {
			case <-ctx.Done():
				log.Println("[JOB] Stopping scheduler...")
				return

			case <-dataTicker.C:
				log.Printf("[MEDIA] Syncing media...")
				go syncFn()

			case <-tokenTicker.C:
				log.Printf("[TOKEN] Checking token refresh...")
				go refTok()
			}
		}
	}()
}
