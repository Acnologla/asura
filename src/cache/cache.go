package cache

import (
	"sync"
	"time"

	"github.com/andersfylling/disgord"
)

const N = 5

type User struct {
	IsBan     bool
	Timestamp time.Time
}

type Cache struct {
	Values map[disgord.Snowflake]User
	*sync.RWMutex
}

func Get(cache Cache, id disgord.Snowflake) (bool, bool) {
	cache.RLock()
	defer cache.RUnlock()
	user, ok := cache.Values[id]
	return user.IsBan, ok
}

func Set(cache Cache, id disgord.Snowflake, value bool) {
	cache.Lock()
	defer cache.Unlock()
	cache.Values[id] = User{
		IsBan:     value,
		Timestamp: time.Now(),
	}
}

func New() Cache {
	cache := Cache{
		Values:  map[disgord.Snowflake]User{},
		RWMutex: &sync.RWMutex{},
	}
	go func() {
		ticker := time.NewTicker(N * time.Minute)
		for range ticker.C {
			cache.Lock()
			for id, user := range cache.Values {
				if time.Since(user.Timestamp).Minutes() >= N {
					delete(cache.Values, id)
				}
			}
			cache.Unlock()
		}
	}()
	return cache
}
