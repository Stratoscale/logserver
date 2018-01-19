package engine

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type Cache interface {
	Set(key string, resp *Response)
	Get(key string) (*Response, bool)
}

func NewCache(exp time.Duration) Cache {
	return &responseCache{
		cache: cache.New(exp, exp),
	}
}

type responseCache struct {
	cache *cache.Cache
}

func (c *responseCache) Set(key string, resp *Response) {
	c.cache.Set(key, resp, cache.DefaultExpiration)
}

func (c *responseCache) Get(key string) (*Response, bool) {
	resp, ok := c.cache.Get(key)
	if !ok {
		return nil, false
	}
	return resp.(*Response), true
}
