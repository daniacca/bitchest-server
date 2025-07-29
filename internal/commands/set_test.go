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

func TestSetCommand_WithNX_Success(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Set key with NX when it doesn't exist
	out, err := cmd.Execute([]string{"foo", "bar", "NX"}, store)
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

func TestSetCommand_WithNX_Failure(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Set key first
	store.Set("foo", &db.StringValue{Val: "existing"})

	// Try to set with NX when key exists
	out, err := cmd.Execute([]string{"foo", "bar", "NX"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("Expected empty response for NX failure, got %q", out)
	}

	// Check that original value is unchanged
	val, ok := store.Get("foo")
	if !ok {
		t.Fatal("Key should still exist")
	}
	if str, ok := val.(*db.StringValue); !ok || str.Get() != "existing" {
		t.Errorf("Expected 'existing', got %v", val)
	}
}

func TestSetCommand_WithXX_Success(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Set key first
	store.Set("foo", &db.StringValue{Val: "existing"})

	// Set with XX when key exists
	out, err := cmd.Execute([]string{"foo", "bar", "XX"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "+OK") {
		t.Errorf("Expected +OK response, got %q", out)
	}

	val, ok := store.Get("foo")
	if !ok {
		t.Fatal("Key should exist")
	}
	if str, ok := val.(*db.StringValue); !ok || str.Get() != "bar" {
		t.Errorf("Expected 'bar', got %v", val)
	}
}

func TestSetCommand_WithXX_Failure(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Try to set with XX when key doesn't exist
	out, err := cmd.Execute([]string{"foo", "bar", "XX"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if out != "" {
		t.Errorf("Expected empty response for XX failure, got %q", out)
	}

	// Check that key was not set
	_, ok := store.Get("foo")
	if ok {
		t.Fatal("Key should not exist")
	}
}

func TestSetCommand_WithNXAndExpiration(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Set key with NX and EX
	out, err := cmd.Execute([]string{"foo", "bar", "NX", "EX", "10"}, store)
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

func TestSetCommand_WithXXAndExpiration(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Set key first
	store.Set("foo", &db.StringValue{Val: "existing"})

	// Set with XX and EX
	out, err := cmd.Execute([]string{"foo", "bar", "XX", "EX", "10"}, store)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "+OK") {
		t.Errorf("Expected +OK response, got %q", out)
	}

	val, ok := store.Get("foo")
	if !ok {
		t.Fatal("Key should exist")
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

	// Test missing expiration time
	_, err = cmd.Execute([]string{"foo", "bar", "EX"}, store)
	if err == nil {
		t.Fatal("Expected error on missing expiration time")
	}
}

func TestSetCommand_MultipleNXOrXX(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Test multiple NX
	_, err := cmd.Execute([]string{"foo", "bar", "NX", "NX"}, store)
	if err == nil {
		t.Fatal("Expected error on multiple NX")
	}

	// Test multiple XX
	_, err = cmd.Execute([]string{"foo", "bar", "XX", "XX"}, store)
	if err == nil {
		t.Fatal("Expected error on multiple XX")
	}

	// Test NX and XX together
	_, err = cmd.Execute([]string{"foo", "bar", "NX", "XX"}, store)
	if err == nil {
		t.Fatal("Expected error on NX and XX together")
	}

	// Test XX and NX together
	_, err = cmd.Execute([]string{"foo", "bar", "XX", "NX"}, store)
	if err == nil {
		t.Fatal("Expected error on XX and NX together")
	}
}

func TestSetCommand_InvalidOption(t *testing.T) {
	store := db.NewDB()
	cmd := &SetCommand{}

	// Test invalid option
	_, err := cmd.Execute([]string{"foo", "bar", "INVALID"}, store)
	if err == nil {
		t.Fatal("Expected error on invalid option")
	}
}