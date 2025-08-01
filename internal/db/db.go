package db

import (
	"sync"
	"time"
)

type Stats struct {
	Keys int // Number of keys in the database
	MemoryUsage int // Memory usage in bytes - total size of all values in the database
	MemoryPerKey int // Memory usage per key - average size of all values in the database
	PeakMemoryUsage int // Peak memory usage in bytes - maximum memory usage in the database
	NumberOfExpiredKeys int // Number of expired keys in the database
	DataSize int // Number of entries in the database
}

type InMemoryDB struct {
	data map[string]Value
	mu   sync.RWMutex
	stats Stats
}

func NewDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]Value),
		stats: Stats{
			Keys: 0,
			MemoryUsage: 0,
			MemoryPerKey: 0,
			PeakMemoryUsage: 0,
			NumberOfExpiredKeys: 0,
			DataSize: 0,
		},
	}
}

func (db *InMemoryDB) Set(key string, val Value) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Add the key to the database
	db.data[key] = val

	// Update stats
	db.stats.Keys++
	db.stats.DataSize++
	db.stats.MemoryUsage += val.Size()
	db.stats.MemoryPerKey = db.stats.MemoryUsage / db.stats.DataSize
	db.stats.PeakMemoryUsage = max(db.stats.PeakMemoryUsage, db.stats.MemoryUsage)
}

func (db *InMemoryDB) Get(key string) (Value, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	if !ok {
		return nil, false
	}
	
	if val.IsExpired() {
		// Remove expired key (we need to upgrade to write lock)
		db.mu.RUnlock()
		db.mu.Lock()
		delete(db.data, key)
		db.stats.NumberOfExpiredKeys++
		db.mu.Unlock()
		db.mu.RLock() // Re-acquire read lock, since we're in a defer
		return nil, false
	}
	
	return val, true
}

func (db *InMemoryDB) Keys() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	keys := make([]string, 0, len(db.data))
	for k, v := range db.data {
		if !v.IsExpired() {
			keys = append(keys, k)
		}
	}
	return keys
}

func (db *InMemoryDB) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, exists := db.data[key]; exists {
		// Get the size of the key
		removedSize := db.data[key].Size()

		// Remove the key from the database
		delete(db.data, key)

		// Update stats
		db.stats.Keys--
		db.stats.DataSize--
		db.stats.MemoryUsage -= removedSize
		if db.stats.DataSize > 0 {
			db.stats.MemoryPerKey = db.stats.MemoryUsage / db.stats.DataSize
		} else {
			db.stats.MemoryPerKey = 0
		}
		return true
	}
	return false
}

func (db *InMemoryDB) FlushAll() {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Clear the database
	db.data = make(map[string]Value)

	// Reset stats
	db.stats = Stats{
		Keys: 0,
		MemoryUsage: 0,
		MemoryPerKey: 0,
		PeakMemoryUsage: db.stats.PeakMemoryUsage,
		NumberOfExpiredKeys: db.stats.NumberOfExpiredKeys,
		DataSize: 0,
	}
}

func (db *InMemoryDB) SetExpiration(key string, seconds int) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	val, exists := db.data[key]
	if !exists {
		return false
	}
	
	// Set expiration time
	expireAt := time.Now().Add(time.Duration(seconds) * time.Second)
	
	// Handle different value types
	switch v := val.(type) {
	case *StringValue:
		v.ExpireAt = &expireAt
	default:
		// For now, only StringValue supports expiration
		return false
	}
	
	return true
}

func (db *InMemoryDB) GetTTL(key string) int {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	val, exists := db.data[key]
	if !exists {
		return -2 // Key doesn't exist
	}
	
	if val.IsExpired() {
		return -2 // Key is expired
	}
	
	// Handle different value types
	switch v := val.(type) {
	case *StringValue:
		if v.ExpireAt == nil {
			return -1 // No expiration set
		}
		ttl := int(time.Until(*v.ExpireAt).Seconds())
		if ttl < 0 {
			return -2 // Expired
		}
		return ttl
	default:
		return -1 // No expiration support for this type
	}
}

// CleanupExpired removes all expired keys from the database
func (db *InMemoryDB) CleanupExpired() int {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	removed := 0
	for key, val := range db.data {
		if val.IsExpired() {
			delete(db.data, key)
			removed++
		}
	}
	
	return removed
}

// GetStats returns the current stats of the database
func (db *InMemoryDB) GetStats() Stats {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.stats
}