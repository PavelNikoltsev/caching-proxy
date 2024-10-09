package cache

import (
	"fmt"
	"sync"
	"time"
)

var Cache *CacheStore

func Init() {
	Cache = &CacheStore{Store: make(map[string][]byte)}
}

type CacheStore struct {
	Store map[string][]byte
	sync.RWMutex
}

func (c *CacheStore) Get(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()
	v, ok := c.Store[key]
	return v, ok
}

func (c *CacheStore) Set(key string, value []byte) {
	c.Lock()
	defer c.Unlock()
	c.Store[key] = value
}

func (c *CacheStore) Clear() {
	c.Lock()
	defer c.Unlock()
	c.Store = make(map[string][]byte)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] Cache cleared\n", currentTime)
}

func (c *CacheStore) StartAutoClear(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			c.Clear()
		}
	}()
}
