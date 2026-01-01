package cache

import (
	"sync"
	"testing"

	"backend-service/internal/instagram"
)

func TestSetAndGetMedia(t *testing.T) {
	store := NewStore()

	media := []instagram.Media{
		{ID: "1"},
		{ID: "2"},
	}

	store.SetMedia(media)

	result := store.GetAllMedia()

	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
}

func TestConcurrentMediaAccess(t *testing.T) {
	store := NewStore()
	store.SetMedia([]instagram.Media{{ID: "1"}})

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.GetAllMedia()
		}()
	}

	wg.Wait()
}
