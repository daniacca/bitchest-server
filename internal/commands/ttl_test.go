package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestTTLCommand(t *testing.T) {
	store := db.NewDB()
	cmd := &TTLCommand{}

	// Test TTL on non-existent key
	result, err := cmd.Execute([]string{"nonexistent"}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != protocol.Integer(-2) {
		t.Errorf("Expected -2, got %s", result)
	}

	// Set a key without expiration
	store.Set("test", &db.StringValue{Val: "value"})
	result, err = cmd.Execute([]string{"test"}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != protocol.Integer(-1) {
		t.Errorf("Expected -1, got %s", result)
	}

	// Set a key with expiration
	store.SetExpiration("test", 10)
	_, err = cmd.Execute([]string{"test"}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	ttl := store.GetTTL("test")
	if ttl < 0 || ttl > 10 {
		t.Errorf("Expected TTL between 0-10, got %d", ttl)
	}

	// Test invalid arguments
	_, err = cmd.Execute([]string{}, store)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}

	_, err = cmd.Execute([]string{"key1", "key2"}, store)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}
} 