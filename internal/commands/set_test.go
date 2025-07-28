package commands

import (
	"strings"
	"testing"

	"github.com/daniacca/bitchest/internal/db"
)

func TestSetCommand_Valid(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	out, err := cmd.Execute([]string{"foo", "bar"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "+OK") {
		t.Errorf("Expected +OK response, got %q", out)
	}

	val, ok := store.Get("foo")
	if !ok {
		t.Fatal("Key not set in store")
	}
	if str, ok := val.(*db.StringValue); !ok || str.Get() != "bar" {
		t.Errorf("Expected 'bar', got %v", val)
	}
}

func TestSetCommand_WithExpiration(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	out, err := cmd.Execute([]string{"foo", "bar", "EX", "10"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "+OK") {
		t.Errorf("Expected +OK response, got %q", out)
	}

	val, ok := store.Get("foo")
	if !ok {
		t.Fatal("Key not set in store")
	}
	if str, ok := val.(*db.StringValue); !ok || str.Get() != "bar" {
		t.Errorf("Expected 'bar', got %v", val)
	}

	// Check that expiration is set
	if str, ok := val.(*db.StringValue); ok {
		if str.ExpireAt == nil {
			t.Error("Expected expiration to be set")
		}
		ttl := store.GetTTL("foo")
		if ttl < 0 || ttl > 10 {
			t.Errorf("Expected TTL between 0-10, got %d", ttl)
		}
	}
}

func TestSetCommand_NoValueProvided(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	_, err := cmd.Execute([]string{"foo"}, store)
	if err == nil {
		t.Fatal("Expected error on wrong number of args")
	}
}

func TestSetCommand_InvalidExpiration(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Test invalid expiration time
	_, err := cmd.Execute([]string{"foo", "bar", "EX", "invalid"}, store)
	if err == nil {
		t.Fatal("Expected error on invalid expiration time")
	}

	// Test negative expiration time
	_, err = cmd.Execute([]string{"foo", "bar", "EX", "-1"}, store)
	if err == nil {
		t.Fatal("Expected error on negative expiration time")
	}
}