package cache

import (
	"log"
	"sync"

	"backend-service/internal/instagram"
)

type Store struct {
	mu    sync.RWMutex
	media map[string]instagram.Media
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

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.media = make(map[string]instagram.Media)
}
