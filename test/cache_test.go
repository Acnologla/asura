package test

import (
	"asura/src/cache"
	"testing"
)

func TestCache(t *testing.T) {
	_cache := cache.New()
	t.Run("TestCacheSet", func(t *testing.T) {
		cache.Set(_cache, 0, true)
	})
	t.Run("TestCacheGet", func(t *testing.T) {
		result, _ := cache.Get(_cache, 0)
		if !result {
			t.Errorf("This must return true")
		}
		result, _ = cache.Get(_cache, 1)
		if result {
			t.Errorf("This must return false")
		}
	})
}
