package db

import (
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	db := NewDB()

	db.Set("foo", &StringValue{Val: "bar"})

	val, ok := db.Get("foo")
	if !ok {
		t.Fatal("Expected key to exist")
	}

	strVal, ok := val.(*StringValue)
	if !ok {
		t.Fatal("Expected value to be of type StringValue")
	}

	if strVal.Get() != "bar" {
		t.Fatalf("Expected 'bar', got '%s'", strVal.Get())
	}
}

func TestDelete(t *testing.T) {
	db := NewDB()

	db.Set("foo", &StringValue{Val: "bar"})
	deleted := db.Delete("foo")
	if !deleted {
		t.Fatal("Expected key to be deleted")
	}

	_, ok := db.Get("foo")
	if ok {
		t.Fatal("Expected key to be gone")
	}
}

func TestFlushAll(t *testing.T) {
	db := NewDB()

	db.Set("a", &StringValue{Val: "1"})
	db.Set("b", &StringValue{Val: "2"})
	db.FlushAll()

	keys := db.Keys()
	if len(keys) != 0 {
		t.Fatalf("Expected 0 keys, got %d", len(keys))
	}
}

func TestExpiration(t *testing.T) {
	db := NewDB()

	// Set a key with expiration
	db.Set("foo", &StringValue{Val: "bar"})
	success := db.SetExpiration("foo", 1)
	if !success {
		t.Fatal("Expected expiration to be set successfully")
	}

	// Check TTL
	ttl := db.GetTTL("foo")
	if ttl < 0 || ttl > 1 {
		t.Fatalf("Expected TTL between 0-1, got %d", ttl)
	}

	// Wait for expiration
	time.Sleep(1100 * time.Millisecond)

	// Key should be expired
	_, ok := db.Get("foo")
	if ok {
		t.Fatal("Expected key to be expired")
	}

	// TTL should return -2 for expired key
	ttl = db.GetTTL("foo")
	if ttl != -2 {
		t.Fatalf("Expected TTL -2 for expired key, got %d", ttl)
	}
}

func TestTTL(t *testing.T) {
	db := NewDB()

	// Test TTL on non-existent key
	ttl := db.GetTTL("nonexistent")
	if ttl != -2 {
		t.Fatalf("Expected TTL -2 for non-existent key, got %d", ttl)
	}

	// Test TTL on key without expiration
	db.Set("foo", &StringValue{Val: "bar"})
	ttl = db.GetTTL("foo")
	if ttl != -1 {
		t.Fatalf("Expected TTL -1 for key without expiration, got %d", ttl)
	}

	// Test TTL on key with expiration
	success := db.SetExpiration("foo", 10)
	if !success {
		t.Fatal("Expected expiration to be set successfully")
	}
	ttl = db.GetTTL("foo")
	if ttl < 0 || ttl > 10 {
		t.Fatalf("Expected TTL between 0-10, got %d", ttl)
	}
}

func TestSetExpirationOnNonExistentKey(t *testing.T) {
	db := NewDB()

	success := db.SetExpiration("nonexistent", 10)
	if success {
		t.Fatal("Expected failure when setting expiration on non-existent key")
	}
}

func TestCleanupExpired(t *testing.T) {
	db := NewDB()

	// Set keys with different expiration times
	db.Set("expired", &StringValue{Val: "expired"})
	db.Set("valid", &StringValue{Val: "valid"})

	// Set expired key to past time
	expiredTime := time.Now().Add(-1 * time.Second)
	if val, ok := db.Get("expired"); ok {
		if strVal, ok := val.(*StringValue); ok {
			strVal.ExpireAt = &expiredTime
		}
	}

	// Set valid key to future time
	validTime := time.Now().Add(10 * time.Second)
	if val, ok := db.Get("valid"); ok {
		if strVal, ok := val.(*StringValue); ok {
			strVal.ExpireAt = &validTime
		}
	}

	// Cleanup expired keys
	removed := db.CleanupExpired()
	if removed != 1 {
		t.Fatalf("Expected 1 expired key to be removed, got %d", removed)
	}

	// Check that only valid key remains
	keys := db.Keys()
	if len(keys) != 1 || keys[0] != "valid" {
		t.Fatalf("Expected only 'valid' key to remain, got %v", keys)
	}
}

func TestGetStats(t *testing.T) {
	db := NewDB()

	db.Set("foo", &StringValue{Val: "bar"})
	stats := db.GetStats()

	if stats.Keys != 1 {
		t.Fatalf("Expected 1 key, got %d", stats.Keys)
	}
}