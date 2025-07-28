package commands

import (
	"testing"

	"github.com/daniacca/bitchest/internal/db"
	"github.com/daniacca/bitchest/internal/protocol"
)

func TestExpireCommand(t *testing.T) {
	store := db.NewDB()
	cmd := &ExpireCommand{}

	// Set a key first
	store.Set("test", &db.StringValue{Val: "value"})

	// Test setting expiration
	result, err := cmd.Execute([]string{"test", "10"}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != protocol.Integer(1) {
		t.Errorf("Expected 1, got %s", result)
	}

	// Test TTL
	ttl := store.GetTTL("test")
	if ttl < 0 || ttl > 10 {
		t.Errorf("Expected TTL between 0-10, got %d", ttl)
	}

	// Test setting expiration on non-existent key
	result, err = cmd.Execute([]string{"nonexistent", "10"}, store)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != protocol.Integer(0) {
		t.Errorf("Expected 0, got %s", result)
	}

	// Test invalid arguments
	_, err = cmd.Execute([]string{"test"}, store)
	if err == nil {
		t.Error("Expected error for wrong number of arguments")
	}

	_, err = cmd.Execute([]string{"test", "invalid"}, store)
	if err == nil {
		t.Error("Expected error for invalid expiration time")
	}

	_, err = cmd.Execute([]string{"test", "-1"}, store)
	if err == nil {
		t.Error("Expected error for negative expiration time")
	}
} 