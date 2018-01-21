package cache

import (
	"time"

	"github.com/bluele/gcache"
)

const (
	defaultExpiration = time.Minute * 5
	defaultSize       = 100
)

// Config are cache configuration
type Config struct {
	Expiration time.Duration `json:"expiration"`
	Size       int           `json:"size"`
}

func New(c Config) gcache.Cache {
	if c.Expiration == 0 {
		c.Expiration = defaultExpiration
	}
	if c.Size == 0 {
		c.Size = defaultSize
	}
	return gcache.New(c.Size).Expiration(c.Expiration).Build()
}
