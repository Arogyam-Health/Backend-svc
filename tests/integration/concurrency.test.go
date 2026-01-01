package integration

import (
	"sync"
	"testing"

	"backend-service/internal/cache"
	"backend-service/internal/instagram"
)

func TestConcurrentMediaAccess(t *testing.T) {
	store := cache.NewStore()
	store.SetMedia([]instagram.Media{{ID: "1"}})

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.GetAllMedia()
		}()
	}

	wg.Wait()
}
