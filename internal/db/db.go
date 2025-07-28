package db

import "sync"

type InMemoryDB struct {
	data map[string]Value
	mu   sync.RWMutex
}

func NewDB() *InMemoryDB {
	return &InMemoryDB{
		data: make(map[string]Value),
	}
}

func (db *InMemoryDB) Set(key string, val Value) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[key] = val
}

func (db *InMemoryDB) Get(key string) (Value, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	val, ok := db.data[key]
	return val, ok
}

func (db *InMemoryDB) Keys() []string {
	db.mu.RLock()
	defer db.mu.RUnlock()
	keys := make([]string, 0, len(db.data))
	for k := range db.data {
		keys = append(keys, k)
	}
	return keys
}

func (db *InMemoryDB) Delete(key string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	if _, exists := db.data[key]; exists {
		delete(db.data, key)
		return true
	}
	return false
}

func (db *InMemoryDB) FlushAll() {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data = make(map[string]Value)
}