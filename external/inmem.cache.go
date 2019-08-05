package external

import (
	"sync"

	"github.com/NickRI/fiatconv/converter/domain/entries"
)

type InMemCache struct {
	cache map[string]entries.CachedAmount
	sync.RWMutex
}

func NewInMemCache() KeyValStorage {
	return &InMemCache{cache: make(map[string]entries.CachedAmount)}
}

func (i *InMemCache) Get(key string) (entries.CachedAmount, error) {
	i.RLock()
	ca, ok := i.cache[key]
	i.RUnlock()
	if !ok {
		return ca, entries.KeyNotFoundError
	}
	return ca, nil
}

func (i *InMemCache) Set(key string, val entries.CachedAmount) error {
	i.Lock()
	i.cache[key] = val
	i.Unlock()
	return nil
}
