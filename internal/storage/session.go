package storage

import (
	"context"
	app "github.com/VMAnalytic/alarm-bot/internal"
	"sync"
	"time"
)

const ttl = 5 * time.Minute

type MemorySessionStorage struct {
	sync.RWMutex
	sessions          map[int]Item
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
}

func NewMemorySessionStorage() *MemorySessionStorage {
	items := make(map[int]Item)
	storage := &MemorySessionStorage{
		sessions:          items,
		defaultExpiration: ttl,
		cleanupInterval:   ttl,
	}
	storage.StartGC()

	return storage
}

type Item struct {
	Value      app.Session
	Expiration int64
	Created    time.Time
}

func (s *MemorySessionStorage) Add(ctx context.Context, session *app.Session) {
	var expiration int64

	if session == nil {
		return
	}
	expiration = time.Now().Add(s.defaultExpiration).UnixNano()
	s.Lock()

	defer s.Unlock()

	s.sessions[session.ID] = Item{
		Value:      *session,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (s *MemorySessionStorage) ExistInState(ctx context.Context, ID int, state app.SessionState) bool {
	s.RLock()

	defer s.RUnlock()

	item, found := s.sessions[ID]

	// cache not found
	if !found {
		return false
	}

	if item.Expiration > 0 {
		// cache expired
		if time.Now().UnixNano() > item.Expiration {
			return false
		}
	}

	return item.Value.State == state
}

func (s *MemorySessionStorage) Delete(ctx context.Context, ID int) {
	s.Lock()

	defer s.Unlock()

	if _, found := s.sessions[ID]; !found {
		return
	}

	delete(s.sessions, ID)
}

func (s *MemorySessionStorage) StartGC() {
	go s.gc()
}

func (s *MemorySessionStorage) gc() {
	for {
		<-time.After(s.cleanupInterval)

		if s.sessions == nil {
			return
		}

		if keys := s.expiredKeys(); len(keys) != 0 {
			s.clearItems(keys)
		}
	}
}

func (s *MemorySessionStorage) expiredKeys() (keys []int) {
	s.RLock()

	defer s.RUnlock()

	for k, i := range s.sessions {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (s *MemorySessionStorage) clearItems(keys []int) {
	s.Lock()

	defer s.Unlock()

	for _, k := range keys {
		delete(s.sessions, k)
	}
}
