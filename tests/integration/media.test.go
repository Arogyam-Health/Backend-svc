package integration

import (
	"testing"

	"backend-service/internal/cache"
	"backend-service/internal/instagram"
	"backend-service/tests/dummy"
	"backend-service/tests/helpers"
)

func TestMediaBootstrap(t *testing.T) {
	srv := dummy.StartDummyServer()
	defer srv.Close()

	store := cache.NewStore()
	service := helpers.NewTestService("dummy_user")

	media, err := service.FetchMedia()
	if err != nil {
		t.Fatal(err)
	}

	store.SetMedia(media)
	media2 := store.GetAllMedia()

	if len(media2) == 0 {
		t.Fatal("media not cached")
	}
}

func TestMediaSyncAddsNewItems(t *testing.T) {
	store := cache.NewStore()

	store.SetMedia([]instagram.Media{{ID: "1"}})

	newMedia := []instagram.Media{
		{ID: "1"},
		{ID: "2"},
	}

	store.SetMedia(newMedia)

	media := store.GetAllMedia()
	if len(media) != 2 {
		t.Fatal("media sync did not add new items")
	}
}
