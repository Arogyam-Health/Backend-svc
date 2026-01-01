package integration

import (
	"testing"

	"backend-service/internal/cache"
	"backend-service/internal/instagram"
	"backend-service/tests/helpers"
)

func TestInstagramDownDoesNotCrashService(t *testing.T) {
	store := cache.NewStore()

	// pre-fill cache (previous sync)
	store.SetMedia([]instagram.Media{
		{ID: "cached_media"},
	})

	service := helpers.NewTestService("dummy")

	// no dummy server running → IG is DOWN
	_, err := service.FetchMedia()

	if err == nil {
		t.Fatal("expected error when Instagram is down")
	}

	media := store.GetAllMedia()

	// cache must still be intact
	if len(media) == 0 {
		t.Fatal("cached media should not be cleared")
	}
}
