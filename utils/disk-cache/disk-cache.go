package disk_cache

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"sync"
)

// DiskCache represents a simple disk-based cache
type DiskCache struct {
	cacheFile string
	data      map[string]string
	mu        sync.RWMutex
}

// NewDiskCache returns a new DiskCache instance
func NewDiskCache(cacheFile string) (*DiskCache, error) {
	cache := &DiskCache{
		cacheFile: cacheFile,
		data:      make(map[string]string),
	}

	err := cache.loadCache()
	if err != nil {
		return nil, err
	}

	return cache, nil
}

// Set sets a value in the cache
func (c *DiskCache) Set(key, value string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	return c.saveCache()
}

// Get retrieves a value from the cache
func (c *DiskCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.data[key]
	return value, ok
}

func (c *DiskCache) loadCache() error {
	data, err := os.ReadFile(c.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // file does not exist, ignore
		}
		return errors.Wrap(err, "failed to read cache file")
	}

	err = json.Unmarshal(data, &c.data)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal cache data")
	}

	return nil
}

func (c *DiskCache) saveCache() error {
	data, err := json.Marshal(c.data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cache data")
	}

	err = os.WriteFile(c.cacheFile, data, 0644)
	if err != nil {
		return errors.Wrap(err, "failed to write cache file")
	}

	return nil
}
