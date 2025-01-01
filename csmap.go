// Package csmap provides a concurrent-safe hashmap implementation using a sharded architecture.
// The CSMap structure allows for thread-safe access and modification of key-value pairs by partitioning
// the map into multiple shards. Each shard has its own mutex for synchronization, enabling simultaneous
// read and write operations across different shards. The package includes methods for setting, getting,
// and deleting key-value pairs, as well as a mechanism to determine the appropriate shard for a given key
// based on its hashed value. Users can create a new CSMap instance by specifying the number of shards,
// optimizing performance in multi-threaded environments.
package csmap

import (
	"sync"
	"unsafe"
)

// CSMap is a concurrent hashmap structure that holds shards for thread-safe access.
type CSMap[K comparable, V any] struct {
	shards []*Shard[K, V] // Array of shards that partition the map
	length int            // Number of shards
}

// NewCSMap creates a new CSMap with the specified number of shards.
func NewCSMap[K comparable, V any](length int) *CSMap[K, V] {
	shards := make([]*Shard[K, V], length)
	for i := 0; i < length; i++ {
		shards[i] = &Shard[K, V]{
			m: make(map[K]V),
		}
	}
	return &CSMap[K, V]{
		shards: shards,
		length: length,
	}
}

// Set adds or updates the value for a given key in the map.
func (c *CSMap[K, V]) Set(key K, value V) {
	s := c.getShard(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

// Get retrieves the value associated with the given key.
// It returns the value and a boolean indicating if the key exists.
func (c *CSMap[K, V]) Get(key K) (V, bool) {
	s := c.getShard(key)
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

// Delete removes the key and its associated value from the map.
func (c *CSMap[K, V]) Delete(key K) {
	s := c.getShard(key)
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
}

// getShard returns the shard that should contain the given key.
func (c *CSMap[K, V]) getShard(key K) *Shard[K, V] {
	return c.shards[c.hash(key)] // Retrieve the appropriate shard based on the hash of the key
}

// hash computes a hash for the key to determine which shard it belongs to.
func (c *CSMap[K, V]) hash(key K) uintptr {
	// Use unsafe pointer conversion to get an uintptr representation of the key
	return *(*uintptr)(unsafe.Pointer(&key)) % uintptr(c.length)
}

// Shard is a structure that holds a portion of the map and provides synchronization.
type Shard[K comparable, V any] struct {
	m  map[K]V      // The map storing key-value pairs
	mu sync.RWMutex // Mutex for synchronizing access to the map
}
