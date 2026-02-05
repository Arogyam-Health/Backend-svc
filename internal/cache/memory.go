package cache

import (
	"log"
	"sort"
	"sync"
	"time"

	"backend-service/internal/instagram"
)

type Store struct {
	mu        sync.RWMutex
	media     map[string]instagram.Media
	updatedAt time.Time
}

func NewStore() *Store {
	return &Store{
		media: make(map[string]instagram.Media),
	}
}

func (s *Store) SetMedia(list []instagram.Media) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, media := range list {
		s.media[media.ID] = media
	}
	s.updatedAt = time.Now()
	log.Printf("[CACHE] Updated %d media items at %v", len(list), s.updatedAt.Format(time.RFC3339))
}

func (s *Store) GetByIDs(ids []string) []instagram.Media {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log.Printf("Fetching media by IDs: %v", ids)

	result := make([]instagram.Media, 0, len(ids))
	for _, id := range ids {
		if media, exists := s.media[id]; exists {
			result = append(result, media)
		}
	}
	return result
}

func (s *Store) GetAllMedia() []instagram.Media {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]instagram.Media, 0, len(s.media))
	for _, media := range s.media {
		result = append(result, media)
	}
	return result
}

func (s *Store) GetAllMediaIDs(limit int, mediaType string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Collect media that matches the filter
	filtered := make([]instagram.Media, 0, len(s.media))
	for _, media := range s.media {
		if mediaType != "" && media.MediaType != mediaType {
			continue
		}
		filtered = append(filtered, media)
	}

	// Sort by timestamp descending (latest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp > filtered[j].Timestamp
	})

	// Extract IDs with limit
	result := make([]string, 0, len(filtered))
	for _, media := range filtered {
		result = append(result, media.ID)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.media = make(map[string]instagram.Media)
	s.updatedAt = time.Time{}
}

// IsFresh returns true if cache was updated within the last hour
func (s *Store) IsFresh() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.updatedAt.IsZero() {
		return false
	}
	return time.Since(s.updatedAt) < time.Hour
}

// GetLastUpdateTime returns when the cache was last updated
func (s *Store) GetLastUpdateTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}

// HasMedia checks if specific media IDs exist in cache
func (s *Store) HasMedia(ids []string) (bool, []string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	missing := []string{}
	for _, id := range ids {
		if _, exists := s.media[id]; !exists {
			missing = append(missing, id)
		}
	}
	return len(missing) == 0, missing
}
