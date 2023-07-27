package inmemorycache

import (
	"sync"
	"time"
)

type inmemoryCacheRepository[T any, K comparable] struct {
	m          sync.RWMutex
	cache      map[K]T
	lifeTime   map[K]time.Time
	expiration time.Duration
}

func NewInmemoryCacheRepository[T any, K comparable](expiration time.Duration) *inmemoryCacheRepository[T, K] {
	return &inmemoryCacheRepository[T, K]{
		m:          sync.RWMutex{},
		cache:      make(map[K]T),
		lifeTime:   make(map[K]time.Time),
		expiration: expiration,
	}
}

func (r *inmemoryCacheRepository[T, K]) Add(id K, data T) {
	r.m.Lock()
	r.cache[id] = data
	r.lifeTime[id] = time.Now().Add(r.expiration)
	r.m.Unlock()
}

func (r *inmemoryCacheRepository[T, K]) Get(id K) (t T, exists bool) {
	r.m.RLock()
	defer r.m.RUnlock()

	data, exists := r.cache[id]
	if exists {
		lifeTime, existsLifetime := r.lifeTime[id]
		if existsLifetime {
			if lifeTime.After(time.Now()) {
				return data, true
			} else {
				exists = false
				return
			}
		}

		return data, true
	}

	exists = false
	return
}

func (r *inmemoryCacheRepository[T, K]) Remove(id K) {
	r.m.Lock()
	delete(r.cache, id)
	delete(r.lifeTime, id)
	r.m.Unlock()
}

func (r *inmemoryCacheRepository[T, K]) ClearExpired() {
	now := time.Now()
	r.m.Lock()
	for id, t := range r.lifeTime {
		if t.Before(now) {
			delete(r.cache, id)
			delete(r.lifeTime, id)
		}
	}
	r.m.Unlock()
}
